// Rhiza monitor — the only client script. Server-rendered HTML; this keeps it live.
// A section is an element with data-rz-fragment (a URL) and data-rz-events (event
// kinds, space separated). The script subscribes to /events; when an event of a kind
// a section names arrives, it re-fetches that section's fragment and swaps its
// content. Other sections are left alone.
// On drop: mark the page stale (CSS stops every pulse), show the banner, retry with backoff.
(() => {
  const app = document.querySelector('.rz-app');
  const live = document.querySelector('.rz-live');
  const banner = document.getElementById('rz-stale');
  const kinds = ['change', 'run', 'log', 'inbox', 'daemon'];
  let lastEvent = Date.now(), retry = 1000, es;

  const rel = (t) => {
    const s = Math.round((Date.now() - t) / 1000);
    if (s < 60) return s + 's ago';
    if (s < 3600) return Math.round(s / 60) + 'm ago';
    return Math.round(s / 3600) + 'h ago';
  };
  const setStale = (on) => {
    if (app) app.classList.toggle('is-stale', on);
    if (live) {
      live.classList.toggle('rz-live--stale', on);
      live.lastChild.textContent = on ? 'stale' : 'live';
    }
    if (banner) banner.hidden = !on;
  };
  const refresh = async (el) => {
    try {
      const r = await fetch(el.dataset.rzFragment, { cache: 'no-store' });
      if (r.ok) el.innerHTML = await r.text();
    } catch (_) { /* the stream's own error marks the page stale */ }
  };
  const onEvent = (e) => {
    lastEvent = Date.now();
    document.querySelectorAll('[data-rz-fragment][data-rz-events]').forEach((el) => {
      if (el.dataset.rzEvents.split(/\s+/).includes(e.type)) refresh(el);
    });
  };
  const connect = () => {
    es = new EventSource(((document.body.dataset.rzBase) || '') + '/events');
    es.onopen = () => {
      lastEvent = Date.now();
      // events missed while down: bring every section up to date
      if (retry > 1000) document.querySelectorAll('[data-rz-fragment]').forEach(refresh);
      retry = 1000; setStale(false);
    };
    kinds.forEach((k) => es.addEventListener(k, onEvent));
    es.onerror = () => {
      es.close(); setStale(true);
      setTimeout(connect, retry); retry = Math.min(retry * 2, 30000);
    };
  };
  // relative times: <time datetime="..."> re-rendered every 15s
  setInterval(() => {
    document.querySelectorAll('time[data-rel]').forEach((t) => { t.textContent = rel(Date.parse(t.dateTime)); });
    const since = document.getElementById('rz-stale-since');
    if (since) since.textContent = rel(lastEvent);
  }, 15000);
  // project picker: the header button loads the sheet from /picker for this page
  const pickBtn = document.querySelector('[data-rz-picker]');
  if (pickBtn) {
    let host = null;
    const close = () => { if (host) host.remove(); host = null; pickBtn.setAttribute('aria-expanded', 'false'); };
    pickBtn.addEventListener('click', async () => {
      if (host) return close();
      try {
        const r = await fetch(pickBtn.dataset.rzPicker + '?path=' + encodeURIComponent(location.pathname), { cache: 'no-store' });
        if (!r.ok) return;
        host = document.createElement('div');
        host.className = 'rz-picker-host';
        host.innerHTML = await r.text();
        document.body.appendChild(host);
        pickBtn.setAttribute('aria-expanded', 'true');
        host.querySelector('[data-rz-picker-close]').addEventListener('click', close);
        const q = host.querySelector('[data-rz-picker-search]');
        if (q) q.addEventListener('input', () => {
          const t = q.value.trim().toLowerCase();
          host.querySelectorAll('.rz-picker-row[data-name]').forEach((a) => { a.hidden = !a.dataset.name.toLowerCase().includes(t); });
        });
      } catch (_) { /* the stream's own error marks the page stale */ }
    });
    document.addEventListener('keydown', (e) => { if (e.key === 'Escape') close(); });
  }
  connect();
})();
