// vi: ft=html

package status

const homeTmpl = `
<!doctype html>
<html>
<head>
  <link href="data:image/x-icon;base64,AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAD/hACb/4QAm/+EAJv/hACb+sCCm/rAgpv6wIKb+sCCm/rAgpv6wIKb/966m//eupv/3rqb/966m//eupv/3rqb/4QAm/+EAJv/hACb/4QAm/rAgpv6wIKb+sCCm/rAgpv6wIKb+sCCm//eupv/3rqb/966m//eupv/3rqb/966m/+EAJv/hACb/4QAm/+EAJv6wIKb+sCCm/rAgpv6wIKb+sCCm/rAgpv/3rqb/966m//eupv/3rqb/966m//eupv/hACb/4QAm/+EAJv/hACb+sCCm/rAgpv6wIKb+sCCm/rAgpv6wIKb/966m//eupv/3rqb/966m//eupv/3rqb/4QAm/+EAJv/hACb/4QAm/rAgpv6wIKb+sCCm/rAgpv6wIKb+sCCm//eupv/3rqb/966m//eupv/3rqb/966m/+EAJv/hACb/4QAm/+EAJv6wIKb+sCCm/rAgpv6wIKb+sCCm/rAgpv/3rqb/966m//eupv/3rqb/966m//eupv/hACb/4QAm/+EAJv/hACb+sCCm/rAgpv6wIKb+sCCm/rAgpv6wIKb/966m//eupv/3rqb/966m//eupv/3rqbAAAAAAAAAAAAAAAAAAAAAPrAgpv6wIKb+sCCm/rAgpv6wIKb+sCCm//eupv/3rqb/966m//eupv/3rqb/966mwAAAAAAAAAAAAAAAAAAAAD6wIKb+sCCm/rAgpv6wIKb+sCCm/rAgpv/3rqb/966m//eupv/3rqb/966m//eupsAAAAAAAAAAAAAAAAAAAAA+sCCm/rAgpv6wIKb+sCCm/rAgpv6wIKb/966m//eupv/3rqb/966m//eupv/3rqbAAAAAAAAAAAAAAAAAAAAAPrAgpv6wIKb+sCCm/rAgpv6wIKb+sCCm//eupv/3rqb/966m//eupv/3rqb/966mwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/3rqb/966m//eupv/3rqb/966m//eupsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA/966m//eupv/3rqb/966m//eupv/3rqbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP/eupv/3rqb/966m//eupv/3rqb/966mwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/3rqb/966m//eupv/3rqb/966m//eupsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA/966m//eupv/3rqb/966m//eupv/3rqbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPAAAADwAAAA8AAAAPAAAAD/wAAA/8AAAP/AAAD/wAAA/8AAAA==" rel="icon" type="image/x-icon" />
  <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">
  <meta http-equiv="refresh" content="10" />
  <!-- TODO: bundle this -->
  <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.3/Chart.bundle.min.js"></script>
  <title>{{ .PageTitle }}</title>
  <style>
    :root {
      --group-header-width: 18px;
      --group-header-colour: #ddd;
      --section-radius: 4px;
      --colour-main: #44475A;
    }
    * {
      margin: 0;
      padding: 0;
    }
    body {
      margin: 0 auto;
      max-width: 500px;
      font-family: monospace;
      background-color: #282A36;
    }
    a {
      color: unset;
      text-decoration: none;
    }
    a:hover {
      text-decoration: underline;
    }
    p {
      font-style: italic;
    }
    td,div {
      border-radius: 5px;
    }
    .right {
      text-align: right;
    }
    .light {
      opacity: 0.3;
    }
    .red {
      color: red;
    }
    .green {
      color: #388E3C;
    }
    .stat-table {
      margin-left: auto;
      text-align: right;
    }
    .stat-table tr td:last-child {
      font-weight: bold;
    }
    .stat-table tr td:last-child::before {
      content: '\00a0';
    }
    .group {
      margin: 12px 0;
    }
    .group.group-show .group-title {
      border: 1px #ccc;
      background-color: var(--colour-main);
      filter: brightness(125%);
      color: #444;
      border-bottom: 1px solid #888;
    }
    .group-items {
      background-color: var(--colour-main);
    }
    .group-title,
    .group-items {
      padding: 0 16px;
    }
    .project {
      display: flex;
      flex-wrap: wrap;
    }
    .project~.project {
      border-top: 1px solid white;
    }
  </style>
</head>
{{ $rootProjects := .Projects }}
{{ $rootGroups := .Groups }}
<body>
  <div class="group">
    <div class="group-title"></div>
    <div class="group-items">
      <div class="project">
        {{ if not (eq .PageTitle "") }}
        <strong>{{ .PageTitle }}</strong>
        {{ end }}
        <table class="stat-table">
          <tr>
            <td>cpu</td>
            <td>{{ printf "%.2f" .Stats.CPU }}% {{ printf "%.0f" .Stats.CPUTemp }}&deg;C</td>
          </tr>
          <tr>
            <td>memory</td>
            <td>{{ .Stats.MemUsed | humanBytes }} / {{ .Stats.MemTotal | humanBytes }}</td>
          </tr>
          <tr>
            <td>load</td>
            <td>{{ .Stats.Load1 }} {{ .Stats.Load5 }} {{ .Stats.Load15 }}</td>
          </tr>
          <tr>
            <td>uptime</td>
            <td>{{ .Stats.Uptime | humanDuration }}</td>
          </tr>
        </table>
      </div>
    </div>
  </div>
  <div class="group" id="chart-group">
    <div class="group-title"></div>
    <div class="group-items">
      <canvas id="chart" height="100"></canvas>
    </div>
  </div>
  {{ range $group, $projects := $rootGroups }}
  <div class="group {{ if not (eq $group "~") }}group-show{{ end }}">
    <div class="group-title">
      {{ if not (eq $group "~") }}<h4>{{ $group }}</h4>{{ end }}
    </div>
    <div class="group-items">
      {{ range $projectName := $projects }}
      <div class="project">
        <p><strong>{{ $projectName }}</strong></p>
        <table class="stat-table aligned-stat-table">
          {{ $project := index $rootProjects $projectName }}
          {{ range $container := $project }}
          <tr class="green">
            {{ if not (eq $container.Link "") }}
              <td><a href="//{{ $container.Link }}" target="_blank">{{ $container.Name }}</a></td>
            {{ else }}
              <td>{{ $container.Name }}</td>
            {{ end }}
            <td>{{ $container.Status }}</td>
          </tr>
          {{ end }}
        </table>
      </div>
      {{ end }}
    </div>
  </div>
  {{ end }}
  {{ if .ShowCredit }}
  <div class="group">
    <div class="group-items light right">
      <i><a target="_blank" href="https://github.com/sentriz/compose-status">compose status</a></i>
    </div>
  </div>
  {{ end }}
  <noscript>
    <style>
      #chart-group {
        display: none;
      }
    </style>
  </noscript>
  <script>
    const elem = document.getElementById('chart');
    const ctx = elem.getContext('2d');
    const cpuDY = {{ js .HistDataCPU }};
    const tempDY = {{ js .HistDataTemp }};
    const dx = [];
    const base = new Date().getTime();
    for (var i = 0; i < {{ js (len .HistDataCPU) }}; i++) {
      dx.unshift(new Date(base - ({{ js .HistPeriod.Milliseconds }} * i)));
    }
    const labelFuncs = [
      // can't use backticks here in go template
      (item, data) => "cpu: " + item.value + "%",
      (item, data) => "temperature: " + item.value + "°C",
    ];
    const chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: dx,
        datasets: [
          {
            data: cpuDY,
            pointRadius: 0,
            fill: false,
            borderColor: 'grey',
            borderWidth: 2,
          },
          {
            data: tempDY,
            pointRadius: 0,
            fill: false,
            borderColor: 'orange',
            borderWidth: 2,
          },
        ]
      },
      options: {
        animation: false,
        legend: {
          display: false
        },
        tooltips: {
          callbacks: {
            label(item, data) {
              return labelFuncs[item.datasetIndex](item, data);
            }
          }
        },
        layout: {},
        scales: {
          xAxes: [{
            type: 'time',
            gridLines: {
              display: false
            },
            ticks: {
              maxRotation: 90,
              minRotation: 90
            },
            time: {
              unit: 'second',
              unitStepSize: 5,
              displayFormats: {
                'second': 'HH:mm',
              }
            }
          }],
          yAxes: [{
            ticks: {
              beginAtZero: true,
              max: 100,
              display: false
            },
            gridLines: {
              drawBorder: false,
            }
          }]
        }
      }
    });
  </script>
</body>
</html>
`
