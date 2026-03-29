package catalog

const HTMLTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Privantix Catalog — Catálogo de datasets</title>
  <style>
    :root {
      --primary: #0b1f3b;
      --accent: #18b7a0;
      --bg: #f8fafc;
      --surface: #ffffff;
      --text: #0f172a;
      --muted: #64748b;
      --border: #e2e8f0;
    }
    * { box-sizing: border-box; }
    body { font-family: 'Segoe UI', system-ui, sans-serif; margin: 0; padding: 24px; background: var(--bg); color: var(--text); }
    .container { max-width: 1400px; margin: 0 auto; }
    h1 { font-size: 2rem; margin: 0 0 8px; color: var(--primary); }
    .subtitle { color: var(--muted); margin-bottom: 24px; }
    
    .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; margin-bottom: 32px; }
    .card { background: var(--surface); border: 1px solid var(--border); border-radius: 12px; padding: 16px; }
    .card-title { color: var(--muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 4px; }
    .card-value { font-size: 1.5rem; font-weight: 700; color: var(--accent); }
    
    .search-box { margin-bottom: 24px; }
    .search-box input { width: 100%; max-width: 400px; padding: 10px 16px; border: 1px solid var(--border); border-radius: 8px; font-size: 1rem; }
    
    .dataset { background: var(--surface); border: 1px solid var(--border); border-radius: 12px; margin-bottom: 16px; overflow: hidden; }
    .dataset-header { padding: 16px 20px; background: #f8fafc; border-bottom: 1px solid var(--border); cursor: pointer; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px; }
    .dataset-header:hover { background: #f1f5f9; }
    .dataset-name { font-weight: 600; color: var(--primary); }
    .dataset-meta { font-size: 0.875rem; color: var(--muted); }
    .badge { display: inline-block; padding: 2px 8px; border-radius: 6px; font-size: 0.75rem; font-weight: 600; background: #e0f2fe; color: #0369a1; }
    .badge-warn { background: #fef3c7; color: #92400e; }
    .dataset-body { padding: 20px; display: none; }
    .dataset-body.open { display: block; }
    
    .cols-table { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
    .cols-table th, .cols-table td { padding: 10px 12px; text-align: left; border-bottom: 1px solid var(--border); }
    .cols-table th { background: #f1f5f9; color: var(--muted); font-weight: 600; }
    .cols-table tr:last-child td { border-bottom: none; }
    
    .rules-list { margin-top: 12px; }
    .rule { padding: 6px 10px; background: #fef3c7; border-radius: 6px; font-size: 0.8rem; margin-bottom: 4px; }
    .gov-flags { display: flex; gap: 8px; margin-top: 12px; flex-wrap: wrap; }
    .gov-flag { font-size: 0.75rem; color: var(--muted); }
  </style>
</head>
<body>
  <div class="container">
    <h1>Privantix Catalog</h1>
    <p class="subtitle">Catálogo de datasets generado desde privantix-inspector. Origen: {{ .SourcePath }}</p>
    
    <div class="summary">
      <div class="card">
        <div class="card-title">Datasets</div>
        <div class="card-value">{{ .TotalDatasets }}</div>
      </div>
      <div class="card">
        <div class="card-title">Columnas totales</div>
        <div class="card-value">{{ .TotalColumns }}</div>
      </div>
      <div class="card">
        <div class="card-title">Reglas detectadas</div>
        <div class="card-value">{{ .TotalRules }}</div>
      </div>
      <div class="card">
        <div class="card-title">Generado</div>
        <div class="card-value" style="font-size: 1rem;">{{ .GeneratedAt.Format "2006-01-02 15:04" }}</div>
      </div>
    </div>

    <div class="search-box">
      <input type="text" id="search" placeholder="Buscar dataset o columna..." onkeyup="filterCatalog()">
    </div>

    <div id="catalog">
      {{ range .Datasets }}
      <div class="dataset" data-name="{{ .Name }}" data-path="{{ .Path }}" data-cols="{{ range .Columns }}{{ .Name }} {{ end }}">
        <div class="dataset-header" onclick="toggleDataset(this)">
          <div>
            <span class="dataset-name">{{ .Name }}</span>
            <div class="dataset-meta">{{ .Path }} · {{ .Extension }} · {{ .RowCountEstimate }} filas · {{ .ColumnCount }} columnas</div>
          </div>
          <div>
            <span class="badge">{{ .Extension }}</span>
            {{ if .GovernanceFlags.HasChecksum }}<span class="badge">checksum</span>{{ end }}
            {{ if .GovernanceFlags.HasACLs }}<span class="badge">ACL</span>{{ end }}
            {{ if gt (len .RulesTriggered) 0 }}<span class="badge badge-warn">{{ len .RulesTriggered }} reglas</span>{{ end }}
          </div>
        </div>
        <div class="dataset-body">
          <table class="cols-table">
            <thead><tr><th>Columna</th><th>Tipo</th><th>Nulos %</th><th>Max length</th></tr></thead>
            <tbody>
              {{ range .Columns }}
              <tr><td><strong>{{ .Name }}</strong></td><td><span class="badge">{{ .InferredType }}</span></td><td>{{ printf "%.2f" .NullPercentage }}</td><td>{{ .MaxLength }}</td></tr>
              {{ end }}
            </tbody>
          </table>
          {{ if gt (len .RulesTriggered) 0 }}
          <div class="rules-list">
            <strong>Reglas:</strong>
            {{ range .RulesTriggered }}
            <div class="rule">{{ .RuleName }}{{ if .Target }} → {{ .Target }}{{ end }} ({{ .Severity }})</div>
            {{ end }}
          </div>
          {{ end }}
          <div class="gov-flags">
            {{ if .GovernanceFlags.HasChecksum }}<span class="gov-flag">✓ Checksum</span>{{ end }}
            {{ if .GovernanceFlags.HasACLs }}<span class="gov-flag">✓ ACL</span>{{ end }}
            {{ if .GovernanceFlags.SamplesHidden }}<span class="gov-flag">Muestras ocultas (gobierno)</span>{{ end }}
          </div>
        </div>
      </div>
      {{ end }}
    </div>
  </div>

  <script>
    function toggleDataset(header) {
      const body = header.nextElementSibling;
      body.classList.toggle('open');
    }
    function filterCatalog() {
      const q = document.getElementById('search').value.toLowerCase();
      document.querySelectorAll('.dataset').forEach(el => {
        const name = (el.dataset.name || '').toLowerCase();
        const path = (el.dataset.path || '').toLowerCase();
        const cols = (el.dataset.cols || '').toLowerCase();
        const match = !q || name.includes(q) || path.includes(q) || cols.includes(q);
        el.style.display = match ? '' : 'none';
      });
    }
  </script>
</body>
</html>`
