// Rhiza monitor — the only client script. Server-rendered HTML; this keeps it live.
// Listens to an SSE stream; each event carries {target: "#id", html: "..."} fragments.
// On drop: mark the page stale (CSS stops every pulse), show the banner, retry with backoff.
(() => {
  const app = document.querySelector('.rz-app');
  const live = document.querySelector('.rz-live');
  const banner = document.getElementById('rz-stale');
  let lastEvent = Date.now(), retry = 1000, es;

  const rel = (t) => {
    const s = Math.round((Date.now() - t) / 1000);
    if (s < 60) return s + 's ago';
    if (s < 3600) return Math.round(s / 60) + 'm ago';
    return Math.round(s / 3600) + 'h ago';
  };
  const setStale = (on) => {
    app.classList.toggle('is-stale', on);
    live.classList.toggle('rz-live--stale', on);
    live.lastChild.textContent = on ? 'stale' : 'live';
    if (banner) banner.hidden = !on;
  };
  const connect = () => {
    es = new EventSource('/events');
    es.onopen = () => { retry = 1000; setStale(false); };
    es.onmessage = (e) => {
      lastEvent = Date.now();
      for (const f of JSON.parse(e.data)) {
        const el = document.querySelector(f.target);
        if (el) el.outerHTML = f.html;
      }
    };
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
  connect();
})();
