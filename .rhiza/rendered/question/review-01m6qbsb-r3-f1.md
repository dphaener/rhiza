# review-01m6qbsb-r3-f1

question · `01M6QCP08T8040EJGG2MS0V4T7`

## statement

Plan review round 3 left finding f1 (contradiction on web-host-setting-shape): The change set's intent says it supersedes only 'the localhost part' of Decision web-ui-in-daemon-read-only, but the plan adds a supersedes link from web-host-setting-shape to that whole decision. That decision also carries the owner's quoted 'First pass will be read only. For monitoring purposes only' and still decides milestone-28 and web-monitor; web-host-setting-shape restates only the host and read-only allowlist, not the monitoring-only scope or its quote, so superseding the whole decision retires owner-quoted provenance the intent keeps. Either supersede with a narrower statement that explicitly carries the unchanged parts (with derived-from to the original decision) or drop the supersedes link and record the localhost change as an amendment. The critic confirmed it: The supersedes link 01M6QCGW5VASCCZX5JR0XDBW7V does carry a note limiting it to the localhost clause, but the model ignores that note: internal/slice/rules.go:203 drops any decision with a live incoming supersedes link. rhiza slice web-monitor --change 01M6QBS no longer lists web-ui-in-daemon-read-only, and web-host-setting-shape does not decide web-monitor, milestone-28, web-api-read-only, web-events-stream, web-templates-embedded, web-live-sections, web-roadmap-page or web-change-set-page. So those capabilities lose the owner-quoted read-only, monitoring-only decision from their slices, which the intent ('The UI stays read-only') keeps. The fix must keep that decision current for them: drop the supersedes link and amend, or supersede with a narrower decision that restates the unchanged parts, derives from the original and decides the same capabilities. How should it be settled?

## Links

### In

- answers [web-host-amends-not-supersedes](../decision/web-host-amends-not-supersedes.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
