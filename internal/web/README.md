# Rhiza Monitor — front end

Plain CSS, one small script, self-hosted fonts and the app icon set, all ready to serve. No build step.

```
site.webmanifest
static/
  css/rhiza-tokens.css   tokens for both themes (theme follows prefers-color-scheme) + @font-face
  css/rhiza.css          every component (.rz-*)
  app.js              SSE live updates, relative times, stale state
  fonts/                 Petrona, Instrument Sans, JetBrains Mono (variable woff2, SIL OFL)
  icons/                 favicon.svg / .ico, PNGs, maskable, head.html
pages/                   static reference HTML for each page, with sample data
```

## Wiring it in

Serve `static/` at `/static/` and `site.webmanifest` at `/site.webmanifest`. Every page's `<head>`:

```html
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" href="/static/icons/favicon.ico" sizes="48x48">
<link rel="icon" href="/static/icons/favicon.svg" type="image/svg+xml">
<link rel="apple-touch-icon" href="/static/icons/apple-touch-icon.png">
<link rel="manifest" href="/site.webmanifest">
<meta name="theme-color" content="#0D0C09" media="(prefers-color-scheme: dark)">
<meta name="theme-color" content="#F2EEE4" media="(prefers-color-scheme: light)">
<link rel="stylesheet" href="/static/css/rhiza-tokens.css">
<link rel="stylesheet" href="/static/css/rhiza.css">
<script src="/static/app.js" defer></script>
```

`<body class="rz">`, then the frame:

```html
<div class="rz-app">
  <header class="rz-header"> wordmark · nav (aria-current="page" on the current link) · daemon pill · live indicator </header>
  <main class="rz-main">
    <div class="rz-stale" id="rz-stale" role="status" hidden> … </div>
    … page content …
  </main>
</div>
```

Copy the exact markup from `pages/*.html`. Dark is the default; light follows the OS. Force one with `<html data-theme="dark">` or `data-theme="light"`.

## Live updates

`app.js` opens an `EventSource` on `/events`. Each message is a JSON array of `{ "target": "#css-selector", "html": "<server-rendered fragment>" }`; each target is replaced in place (`outerHTML`). Give anything that changes an `id` and re-render that fragment on the server.

- `<time datetime="…" data-rel>` elements are re-labelled ("4s ago") every 15s.
- If the stream errors, the page goes stale: `.is-stale` on `.rz-app`, the banner shows, every pulse stops and running states fall back to neutral. It reconnects with backoff (1s → 30s) and clears itself on reconnect.
- Opened from disk, the reference pages have no `/events`, so they show the stale state. That is the intended behaviour.

## Icons

| File | Use |
|---|---|
| `favicon.svg` | Modern browsers' tab icon. The small drawing: heavier strokes, no glow. |
| `favicon.ico` | 16, 32, 48 frames for everything else. |
| `favicon-32.png` | A plain PNG fallback, if you need one. |
| `apple-touch-icon.png` | 180×180, square and opaque; iOS rounds it. |
| `icon-192.png`, `icon-512.png` | Manifest icons, rounded tile. |
| `icon-maskable-512.png`, `icon-maskable.svg` | Full-bleed for Android's adaptive shapes; the r sits inside the 80% safe zone. |
| `icon.svg` | The full drawing, for anything large. |

## Motion and accessibility

Only running things move (a 2.8s breathe on gold dots, a stepped terminal cursor). `prefers-reduced-motion: reduce` stops all of it. Every state pairs colour with a word and a glyph; all text pairs hold at least 4.5:1 in both themes. Links are the only interactive elements; focus is a 2px teal outline.

## Fonts

Petrona, Instrument Sans and JetBrains Mono are under the SIL Open Font License 1.1 (each one's licence text is beside it in `static/fonts/`). Self-hosting and bundling are permitted.
