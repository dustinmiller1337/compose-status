package status

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/oxtoacart/bpool"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type Container struct {
	Name     string
	Status   string
	Link     string
	LastSeen time.Time
	IsDown   bool
	Project  string
}

type Stats struct {
	Load1    float64
	Load5    float64
	Load15   float64
	MemUsed  uint64
	MemTotal uint64
	CPU      float64
}

// we need some sort of unique identifier for containers (when tracking ups
// and downs). the "ID" field from the engine won't do, because we want a
// recreated container with probably a different ID to be considered the same
func (c *Container) ID() string {
	return fmt.Sprintf("%s___%s", c.Project, c.Name)
}

type Controller struct {
	tmpl         *template.Template
	client       *docker.Client
	buffPool     *bpool.BufferPool
	cleanCutoff  time.Duration
	scanInterval time.Duration
	groupLabel   string
	pageTitle    string
	showCredit   bool
	LastProjects map[string]*Container
	LastStats    *Stats
}

func WithCleanCutoff(dur time.Duration) func(*Controller) error {
	return func(c *Controller) error {
		c.cleanCutoff = dur
		return nil
	}
}

func WithScanInternal(dur time.Duration) func(*Controller) error {
	return func(c *Controller) error {
		c.scanInterval = dur
		return nil
	}
}

func WithTitle(title string) func(*Controller) error {
	return func(c *Controller) error {
		c.pageTitle = title
		return nil
	}
}

func WithGroupLabel(label string) func(*Controller) error {
	return func(c *Controller) error {
		c.groupLabel = label
		return nil
	}
}

func WithResume(file []byte) func(*Controller) error {
	return func(c *Controller) error {
		if len(file) <= 0 {
			return nil
		}
		return json.Unmarshal(file, &c.LastProjects)
	}
}

func WithCredit(c *Controller) error {
	c.showCredit = true
	return nil
}

func NewController(options ...func(*Controller) error) (*Controller, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "creating docker client")
	}
	tmpl, err := template.
		New("").
		Funcs(template.FuncMap{
			"humanDate":  humanize.Time,
			"humanBytes": humanize.Bytes,
		}).
		Parse(homeTmpl)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}
	cont := &Controller{
		tmpl:         tmpl,
		client:       client,
		buffPool:     bpool.NewBufferPool(64),
		LastProjects: map[string]*Container{},
		LastStats:    &Stats{},
		// defaults
		cleanCutoff: 3 * 24 * time.Hour,
		pageTitle:   "server status",
		groupLabel:  "com.docker.compose.project",
	}
	for _, option := range options {
		if err := option(cont); err != nil {
			return nil, errors.Wrap(err, "running option")
		}
	}
	return cont, nil
}

func hostFromLabel(label string) string {
	const prefix = "Host:"
	if strings.HasPrefix(label, prefix) {
		trimmed := strings.TrimPrefix(label, prefix)
		return strings.SplitN(trimmed, ",", 2)[0]
	}
	return ""
}

func (c *Controller) GetProjects() error {
	seenIDs := map[string]struct{}{}
	containers, err := c.client.ListContainers(
		docker.ListContainersOptions{},
	)
	if err != nil {
		return errors.Wrap(err, "listing containers")
	}
	// insert the current time for any container we see
	for _, rawTain := range containers {
		project, ok := rawTain.Labels[c.groupLabel]
		if !ok {
			continue
		}
		if len(rawTain.Names) == 0 {
			return fmt.Errorf("%q does not have a name", rawTain.ID)
		}
		tain := &Container{
			Name:     rawTain.Names[0],
			Project:  project,
			Status:   strings.ToLower(rawTain.Status),
			LastSeen: time.Now(),
		}
		if label, ok := rawTain.Labels["traefik.frontend.rule"]; ok {
			tain.Link = hostFromLabel(label)
		}
		seenIDs[tain.ID()] = struct{}{}
		c.LastProjects[tain.ID()] = tain
	}
	// set containers we haven't seen to down, and delete one that haven't
	// seen since since the cutoff
	cutoff := time.Now().Add(-1 * c.cleanCutoff)
	for id, tain := range c.LastProjects {
		if tain.LastSeen.Before(cutoff) {
			delete(c.LastProjects, id)
			continue
		}
		if _, ok := seenIDs[id]; !ok {
			tain.IsDown = true
		}
	}
	return nil
}

func (c *Controller) GetStats() error {
	loadStat, err := load.Avg()
	if err != nil {
		return errors.Wrap(err, "get load stat")
	}
	c.LastStats.Load1 = loadStat.Load1
	c.LastStats.Load5 = loadStat.Load5
	c.LastStats.Load15 = loadStat.Load15
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return errors.Wrap(err, "get mem stat")
	}
	c.LastStats.MemUsed = memStat.Used
	c.LastStats.MemTotal = memStat.Total
	percent, err := cpu.Percent(5*time.Second, false)
	if err != nil {
		return errors.Wrap(err, "get cpu stat")
	}
	if len(percent) != 1 {
		return fmt.Errorf("invalid cpu response")
	}
	c.LastStats.CPU = percent[0]
	return nil
}

func (c *Controller) Start() {
	ticker := time.NewTicker(c.scanInterval)
	for range ticker.C {
		if err := c.GetProjects(); err != nil {
			log.Printf("error getting projects: %v\n", err)
		}
		if err := c.GetStats(); err != nil {
			log.Printf("error getting stats: %v\n", err)
		}
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// group the last seen by project, inserting so that the container
	// names are sorted
	projectMap := map[string][]*Container{}
	for _, tain := range c.LastProjects {
		current := projectMap[tain.Project]
		i := sort.Search(len(current), func(i int) bool {
			return current[i].Name >= tain.Name
		})
		current = append(current, nil)
		copy(current[i+1:], current[i:])
		current[i] = tain
		projectMap[tain.Project] = current
	}
	//
	tmplData := struct {
		PageTitle  string
		ShowCredit bool
		Projects   map[string][]*Container
		Stats      *Stats
	}{
		c.pageTitle,
		c.showCredit,
		projectMap,
		c.LastStats,
	}
	// using a pool of buffers, we can write to one first to catch template
	// errors, which avoids a superfluous write to the response writer
	buff := c.buffPool.Get()
	defer c.buffPool.Put(buff)
	if err := c.tmpl.Execute(buff, tmplData); err != nil {
		http.Error(w, fmt.Sprintf("error executing template: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buff.WriteTo(w)
}
