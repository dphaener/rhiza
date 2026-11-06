# web-host-build-leaves

decision · `01M6QC2N0XFHEYM6VQ8TQKCM0B`

## reasoning

changeLeaves (internal/api/change.go:440) only counts capabilities a change set's ops create that are live and have no live child, and run_build.go:113 and run_review.go build and review from that list, so amending landed leaves in place yields nothing to build. Changing changeLeaves to also count amended capabilities is a tool change outside this change set's intent and an arbiter writes no code. Retiring the two verified leaves and recreating them would drop their links to existing Evidence (web-serve-loopback-verified, web-serve-loopback-qa-1, pr-314-web-ui-server) and contradict Decision web-host-setting-shape, which chose amending over shadowing so one capability owns each behaviour. Adding sibling leaves would make two capabilities claim binding and the status URL. New children under the amended capabilities keep the parents as the full, already-amended description with their evidence, make the parents non-leaves (runs build only childless leaves), and give the build run two leaves scoped to the host delta whose acceptance is lifted from the amended parents, so the work, tests and review match what web-host-setting-shape decided. web-status-shows-host depends-on web-serve-host, since status records the host the listener uses; no depends-on is needed between the amended parents, which are not built.

## statement

Change set 01M6QBS carries its build work on two new childless capabilities it creates, one under each amended landed leaf, so changeLeaves finds them and a build run builds exactly the host delta: web-serve-host (child of web-serve-loopback) owns the rhiza.web.host setting: its default 127.0.0.1, rhiza daemon set validation (literal IPv4/IPv6 via net.ParseIP; empty, hostnames including localhost and values with a port rejected with an error naming the key), the listen address net.JoinHostPort(host, port), and the listener manager moving the listener when only the host changes; its acceptance is the host part of web-serve-loopback's amended acceptance (unit tests for default, 0.0.0.0, ::1 giving [::1]:7787, the rejected values, the host-only move in the fake-listener poll sequence, and the daemon-settings testscript listing rhiza.web.host and accepting 0.0.0.0 and ::1 and rejecting localhost, abc and empty). web-status-shows-host (child of web-url-in-daemon-status, depends-on web-serve-host) owns status recording and showing the effective host as http://<host>:<port>, IPv6 bracketed; its acceptance is the daemon-status testscript steps for 0.0.0.0 and ::1 while the default still reads http://127.0.0.1:7787. The amended parents keep their statements, acceptance and evidence as the whole behaviour and become non-leaf capabilities; no code change to changeLeaves and no retiring of the landed leaves.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b6e0da8b-2](../question/question-b6e0da8b-2.md)
- derived-from [web-host-plan-interview](../source/web-host-plan-interview.md)

### In

- decided-by [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)
- decided-by [web-status-shows-host](../capability/rhiza/web-monitor/web-url-in-daemon-status/web-status-shows-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
