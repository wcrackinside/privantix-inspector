/**
 * Rellena bloques [data-gh-repo="owner/repo"] con datos de
 * GET https://api.github.com/repos/{owner}/{repo}/releases/latest
 */
(function () {
  const GH_ACCEPT = 'application/vnd.github+json';

  function esc(s) {
    if (s == null) return '';
    var d = document.createElement('div');
    d.textContent = s;
    return d.innerHTML;
  }

  function formatDate(iso) {
    if (!iso) return '';
    try {
      return new Date(iso).toLocaleDateString('es-ES', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch (e) {
      return iso;
    }
  }

  function sumDownloads(assets) {
    if (!assets || !assets.length) return 0;
    return assets.reduce(function (n, a) {
      return n + (typeof a.download_count === 'number' ? a.download_count : 0);
    }, 0);
  }

  function fetchLatest(owner, repo) {
    var url =
      'https://api.github.com/repos/' +
      encodeURIComponent(owner) +
      '/' +
      encodeURIComponent(repo) +
      '/releases/latest';
    return fetch(url, { headers: { Accept: GH_ACCEPT }, cache: 'no-store' }).then(function (r) {
      if (r.status === 404) return { _notFound: true };
      if (!r.ok) return { _error: true, status: r.status };
      return r.json();
    });
  }

  function fillBlock(root, data, owner, repo) {
    var meta = root.querySelector('.gh-rel-meta');
    var assetsEl = root.querySelector('.gh-rel-assets');
    var notes = root.querySelector('.gh-rel-notes');
    if (!meta || !assetsEl) return;

    if (data._error) {
      meta.innerHTML =
        '<span class="text-warning">No se pudo cargar la release (HTTP ' +
        esc(data.status) +
        '). <a href="https://github.com/' +
        esc(owner) +
        '/' +
        esc(repo) +
        '/releases" target="_blank" rel="noopener">Ver releases</a></span>';
      assetsEl.innerHTML = '';
      return;
    }
    if (data._notFound) {
      meta.innerHTML =
        'Aún no hay release publicada. <a href="https://github.com/' +
        esc(owner) +
        '/' +
        esc(repo) +
        '/releases" target="_blank" rel="noopener">GitHub Releases</a>';
      assetsEl.innerHTML = '';
      return;
    }

    var total = sumDownloads(data.assets);
    var titlePart = data.name && String(data.name).trim() ? ' · ' + esc(data.name) : '';
    meta.innerHTML =
      '<strong>' +
      esc(data.tag_name) +
      '</strong>' +
      titlePart +
      ' · ' +
      esc(formatDate(data.published_at)) +
      (total > 0 ? ' · <span title="Descargas (assets)">' + total + ' ↓</span>' : '');

    assetsEl.innerHTML = '';
    if (data.assets && data.assets.length) {
      data.assets.forEach(function (a) {
        var dc = typeof a.download_count === 'number' ? ' · ' + a.download_count + ' ↓' : '';
        var link = document.createElement('a');
        link.href = a.browser_download_url;
        link.target = '_blank';
        link.rel = 'noopener';
        link.className = 'btn btn-sm btn-outline-privantix';
        link.innerHTML =
          '<i class="bi bi-download me-1" aria-hidden="true"></i>' + esc(a.name) + esc(dc);
        assetsEl.appendChild(link);
      });
    } else {
      assetsEl.innerHTML =
        '<span class="small" style="color: var(--px-muted);">Sin archivos en esta release.</span>';
    }

    if (notes && data.html_url) {
      notes.href = data.html_url;
      notes.style.display = '';
      notes.textContent = 'Notas de versión';
    }
  }

  async function run() {
    var blocks = document.querySelectorAll('[data-gh-repo]');
    if (!blocks.length) return;

    var byKey = {};
    var keys = [];
    blocks.forEach(function (el) {
      var k = el.getAttribute('data-gh-repo');
      if (k && keys.indexOf(k) === -1) keys.push(k);
    });

    for (var i = 0; i < keys.length; i++) {
      var parts = keys[i].split('/');
      if (parts.length < 2) continue;
      var owner = parts[0];
      var repo = parts.slice(1).join('/');
      byKey[keys[i]] = await fetchLatest(owner, repo);
    }

    blocks.forEach(function (el) {
      var k = el.getAttribute('data-gh-repo');
      var parts = k.split('/');
      var owner = parts[0];
      var repo = parts.slice(1).join('/');
      fillBlock(el, byKey[k] || { _error: true, status: '?' }, owner, repo);
    });

    var hero = document.getElementById('product-hero-release-line');
    if (hero && keys.length) {
      var first = byKey[keys[0]];
      if (first && first.tag_name) {
        hero.textContent =
          'Última release: ' +
          first.tag_name +
          ' · ' +
          formatDate(first.published_at) +
          ' · mismo paquete para todas las herramientas CLI.';
      } else if (first && first._notFound) {
        hero.textContent =
          'Release: pendiente de publicación en GitHub. Las herramientas se distribuyen vía GitHub Releases.';
      } else if (first && first._error) {
        hero.textContent =
          'No se pudo consultar la última versión en GitHub. Revisa la conexión o el límite de la API.';
      }
    }
  }

  run().catch(function () {
    var hero = document.getElementById('product-hero-release-line');
    if (hero) hero.textContent = 'No se pudo cargar la información de releases.';
  });
})();
