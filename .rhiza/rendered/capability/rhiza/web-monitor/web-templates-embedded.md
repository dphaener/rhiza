# web-templates-embedded

capability · `01M6PFJ9A1H6A1YBM55VS6HDQG` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test GETs / and /static/app.js from the handler and gets a rendered page and the script with their content types; the Makefile and CI gain no Node or other build toolchain.

## statement

Pages are server-rendered from html/template files embedded with go:embed, with static assets (CSS, fonts, one small vanilla JS file) served from the binary under /static/, so the web UI needs no build step and no files on disk.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-templates-embedded-verified](../../../evidence/web-templates-embedded-verified.md)
- verifies [web-templates-embedded-qa-1](../../../evidence/web-templates-embedded-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
