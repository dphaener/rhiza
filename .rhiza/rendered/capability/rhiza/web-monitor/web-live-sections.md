# web-live-sections

capability · `01M6PFJ9BEWSJ8QDFFK8NR8N48` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders a sample section with its event names and fragment URL, GETs the fragment URL and gets just that section's HTML; the script is served and names the stale state; no browser test.

## statement

A page section declares the events it listens to and a fragment URL; the embedded script subscribes to /events and re-fetches and swaps only the sections a matching event names, and marks the page stale while the stream is down.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-live-sections-verified](../../../evidence/web-live-sections-verified.md)
- verifies [web-live-sections-qa-1](../../../evidence/web-live-sections-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
