package report

const HTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Privantix Source Inspector Report</title>
  <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
  <style>
    :root {
      --primary: #3b82f6;
      --bg: #f8fafc;
      --surface: #ffffff;
      --text: #0f172a;
      --muted: #64748b;
      --border: #e2e8f0;
    }
    body { font-family: 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 32px; background: var(--bg); color: var(--text); }
    .container { max-width: 1200px; margin: 0 auto; }
    h1, h2, h3 { color: var(--text); margin-top: 0; }
    h1 { font-size: 2.5rem; margin-bottom: 8px; }
    p.muted { color: var(--muted); font-size: 1.1rem; margin-bottom: 32px; }
    
    .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; margin-bottom: 32px; }
    .card { background: var(--surface); border: 1px solid var(--border); border-radius: 12px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
    .card-title { color: var(--muted); font-size: 0.875rem; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 8px; font-weight: 600; }
    .card-value { font-size: 1.8rem; font-weight: 700; color: var(--primary); word-break: break-all; }
    .card-value.path { font-size: 1rem; color: var(--text); }

    .charts-row { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; margin-bottom: 32px; }
    .chart-container { background: var(--surface); border: 1px solid var(--border); border-radius: 12px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }

    .table-container { background: var(--surface); border: 1px solid var(--border); border-radius: 12px; overflow-x: auto; margin-bottom: 32px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
    table { border-collapse: collapse; width: 100%; white-space: nowrap; }
    th, td { border-bottom: 1px solid var(--border); padding: 12px 16px; font-size: 0.875rem; text-align: left; }
    th { background: #f1f5f9; font-weight: 600; color: var(--muted); text-transform: uppercase; letter-spacing: 0.05em; position: sticky; top: 0; }
    tr:last-child td { border-bottom: none; }
    tr:hover { background: #f8fafc; }
    
    .badge { display: inline-block; padding: 2px 8px; border-radius: 9999px; font-size: 0.75rem; font-weight: 600; background: #e0f2fe; color: #0369a1; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Privantix Source Inspector</h1>
    <p class="muted">Portable repository profiling report.</p>
    
    <div class="summary">
      <div class="card">
        <div class="card-title">Target Path</div>
        <div class="card-value path">{{ .Run.Path }}</div>
      </div>
      <div class="card">
        <div class="card-title">Supported Files</div>
        <div class="card-value">{{ .Run.SupportedFiles }}</div>
      </div>
      <div class="card">
        <div class="card-title">Analyzed Files</div>
        <div class="card-value">{{ .Run.AnalyzedFiles }}</div>
      </div>
      <div class="card">
        <div class="card-title">Failed Files</div>
        <div class="card-value">{{ .Run.FailedFiles }}</div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-container">
        <h3>Files by Extension</h3>
        <canvas id="extChart"></canvas>
      </div>
      <div class="chart-container">
        <h3>Estimated Rows per File</h3>
        <canvas id="rowsChart"></canvas>
      </div>
    </div>

    <h2>Files Breakdown</h2>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Path</th><th>Ext</th><th>Owner</th><th>Permissions</th><th>Security (ACLs)</th><th>Checksum</th><th>Encoding</th><th>Delimiter</th><th>Header</th><th>Rows</th><th>Columns</th><th>Rules</th>
          </tr>
        </thead>
        <tbody>
          {{ range .Files }}
          <tr>
            <td>{{ .Path }}</td>
            <td><span class="badge">{{ .Extension }}</span></td>
            <td><span class="muted">{{ .Owner }}</span></td>
            <td><span class="badge">{{ .Permissions }}</span></td>
            <td>
              {{ range .ACLs }}<div style="font-size: 0.75rem; color: #475569;">{{ . }}</div>{{ end }}
            </td>
            <td><code style="font-size: 0.7rem; color: #64748b;">{{ if .Checksum }}{{ .Checksum }}{{ end }}</code></td>
            <td>{{ .Encoding }}</td>
            <td><code>{{ .Delimiter }}</code></td>
            <td>{{ if .HasHeader }}Yes{{ else }}No{{ end }}</td>
            <td>{{ .RowCountEstimate }}</td>
            <td>{{ .ColumnCount }}</td>
            <td>{{ range .RulesTriggered }}{{ .RuleName }}<br>{{ end }}</td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>

    <h2>Column Profiling</h2>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>File</th><th>Column</th><th>Type</th><th>Null %</th><th>Max length</th><th>Samples</th>
          </tr>
        </thead>
        <tbody>
          {{ range .Files }}{{ $filepath := .Path }}{{ range .Columns }}
          <tr>
            <td>{{ $filepath }}</td>
            <td><strong>{{ .Name }}</strong></td>
            <td><span class="badge">{{ .InferredType }}</span></td>
            <td>{{ printf "%.2f" .NullPercentage }}</td>
            <td>{{ .MaxLength }}</td>
            <td><span class="muted">{{ range .SampleValues }}{{ . }}, {{ end }}</span></td>
          </tr>
          {{ end }}{{ end }}
        </tbody>
      </table>
    </div>

    {{ if gt (len .Errors) 0 }}
    <h2>Errors</h2>
    <div class="table-container">
      <table>
        <thead><tr><th>Error Description</th></tr></thead>
        <tbody>
          {{ range .Errors }}<tr><td>{{ . }}</td></tr>{{ end }}
          {{ range .Files }}{{ range .Errors }}<tr><td>{{ . }}</td></tr>{{ end }}{{ end }}
        </tbody>
      </table>
    </div>
    {{ end }}

  </div>

  <script>
    const filesData = [
      {{ range .Files }}
      { ext: "{{ .Extension }}", rows: {{ .RowCountEstimate }}, name: "{{ .Name }}" },
      {{ end }}
    ];

    const extCounts = {};
    filesData.forEach(f => {
      extCounts[f.ext] = (extCounts[f.ext] || 0) + 1;
    });

    const extCtx = document.getElementById('extChart').getContext('2d');
    new Chart(extCtx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(extCounts),
        datasets: [{
          data: Object.values(extCounts),
          backgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#ef4444']
        }]
      },
      options: { responsive: true, maintainAspectRatio: false }
    });

    let sortedFiles = [...filesData].sort((a, b) => b.rows - a.rows);
    if(sortedFiles.length > 20) {
      sortedFiles = sortedFiles.slice(0, 20);
    }
    const rowsCtx = document.getElementById('rowsChart').getContext('2d');
    new Chart(rowsCtx, {
      type: 'bar',
      data: {
        labels: sortedFiles.map(f => f.name.substring(0, 15) + (f.name.length > 15 ? '...' : '')),
        datasets: [{
          label: 'Estimated Rows',
          data: sortedFiles.map(f => f.rows),
          backgroundColor: '#3b82f6'
        }]
      },
      options: { responsive: true, maintainAspectRatio: false }
    });
  </script>
</body>
</html>`
