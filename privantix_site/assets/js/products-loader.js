/**
 * Loads privantix_site/products.json and renders cards into #product-cards-root.
 */
(function () {
  const root = document.getElementById('product-cards-root');
  const errEl = document.getElementById('product-cards-error');
  if (!root) return;

  function esc(s) {
    if (s == null) return '';
    const d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
  }

  fetch('products.json', { cache: 'no-store' })
    .then(function (r) {
      if (!r.ok) throw new Error('No se pudo cargar products.json (' + r.status + ')');
      return r.json();
    })
    .then(function (data) {
      var list = Array.isArray(data) ? data : data.products;
      if (!list || !list.length) throw new Error('products.json vacío o formato inválido');

      root.innerHTML = list
        .map(function (p) {
          var dl = (p.downloads || [])
            .map(function (d) {
              return (
                '<a class="btn btn-sm btn-outline-privantix me-1 mb-1" href="' +
                esc(d.url) +
                '" target="_blank" rel="noopener">' +
                esc(d.label) +
                '</a>'
              );
            })
            .join('');

          var doc = p.documentation
            ? '<a class="btn btn-sm btn-privantix me-1 mb-1" href="' +
              esc(p.documentation) +
              '" target="_blank" rel="noopener">Documentación</a>'
            : '';
          var rm = p.roadmap
            ? '<a class="btn btn-sm btn-outline-secondary me-1 mb-1" href="' +
              esc(p.roadmap) +
              '" target="_blank" rel="noopener">Roadmap</a>'
            : '';

          return (
            '<div class="col-md-6 col-xl-4">' +
            '<div class="feature-card p-4 h-100 d-flex flex-column">' +
            '<div class="d-flex justify-content-between align-items-start gap-2">' +
            '<h5 class="mb-0">' +
            esc(p.name) +
            '</h5>' +
            '<span class="badge rounded-pill badge-soft">' +
            esc(p.version || '—') +
            '</span>' +
            '</div>' +
            '<p class="small mt-3 mb-2 flex-grow-1" style="color: var(--px-muted);">' +
            esc(p.description) +
            '</p>' +
            '<p class="small mb-2 text-white-50"><code>' +
            esc(p.repo || '') +
            '</code></p>' +
            '<div class="mt-auto pt-2">' +
            dl +
            doc +
            rm +
            '</div>' +
            '</div></div>'
          );
        })
        .join('');
    })
    .catch(function (e) {
      if (errEl) {
        errEl.textContent = e.message || String(e);
        errEl.classList.remove('d-none');
      }
    });
})();
