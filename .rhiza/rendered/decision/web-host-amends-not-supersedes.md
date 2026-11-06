# web-host-amends-not-supersedes

decision · `01M6QCQ1JH13B6BXCXJ999GJT0`

## reasoning

The critic confirmed that internal/slice/rules.go:203 drops any decision with a live incoming supersedes link and ignores the link note, so the current link removes the owner-quoted read-only, monitoring-only decision from the slices of milestone-28, web-monitor, web-api-read-only, web-events-stream, web-templates-embedded, web-live-sections, web-roadmap-page and web-change-set-page, which contradicts the intent ('The UI stays read-only'). Of the two fixes the finding offers, dropping the link and amending is smaller and safer: it keeps the owner's exact words on the original decision instead of restating them in an agent-written successor (the model forbids paraphrasing owner decisions), needs no re-linking of eight-plus capabilities, and web-host-setting-shape already decides the leaves whose behaviour changes. The intent's 'supersedes the localhost part' is a partial change, which the model has no partial-supersede for; derived-from with a note expresses it without hiding the rest. A stale 'localhost' word in the original decision is acceptable because the newer decision on the binding leaves is explicit and linked to it.

## statement

Drop the supersedes link 01M6QCGW5VASCCZX5JR0XDBW7V from web-host-setting-shape to web-ui-in-daemon-read-only; web-ui-in-daemon-read-only stays current, owner quote and all, for every capability it decides. The localhost change is recorded as an amendment: web-host-setting-shape links to web-ui-in-daemon-read-only with derived-from, noting that it amends only the 'served on localhost' clause (the UI binds to the address rhiza.web.host names, 127.0.0.1 by default) while the daemon-served, read-only, monitoring-only parts still hold. Its existing decided-by links on m28-read-only, web-serve-loopback, web-url-in-daemon-status and web-serve-host carry the new binding behaviour where it applies. Do not write a narrower superseding decision.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6qbsb-r3-f1](../question/review-01m6qbsb-r3-f1.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
