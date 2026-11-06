# web-api-read-only-verified

evidence · `01M6PSNTH4DZX94GVK3TX5FHH1` · actor [claude-code-cd471593](../actor/claude-code-cd471593.md)

## result

PASS

## statement

Verified at 29d12aca. go test -trimpath -count=1 ./internal/web/... passes: TestAPIAnswersTheHandlersJSON compares GET /api/run.ls, /api/run.ls?limit=5 and /api/roadmap bodies byte-for-byte with canon.Marshal of the api handlers' responses on a fixture repo (init + change.new), the same bytes cmd/rhiza writeResponse prints for --json; TestAPIRefusesWrites shows GET of land, accept, daemon.set, daemon.start, render, ship.merge, add, apply answers 404/405, POST/PUT/DELETE/PATCH answer 405, an unknown query key 400. go test ./internal/api -run TestReadOnlyMarker passes: it asserts ReadOnly on the 16 operations the web UI uses, on no other operation, never with Write or Record, and not on daemon.set, daemon.start, render, ship.merge, the ship/run/session/inbox writes or any log write. Read the diff: Operation.ReadOnly is a new explicit flag separate from Write; ReadOnlyHandler refuses non-GET/HEAD with 405 and any op lacking the marker with 404, builds the request from the query via DecodeRequest, and strips Record/DryRun/Rex/Start/Stdin from a copied Ctx. Spot-read the marked handlers (search, inbox, daemon.status) and found no writes to the log, settings, Rex or GitHub. go vet and gofmt clean; describe-golden testscript passes.

## Links

### Out

- verifies [web-api-read-only](../capability/rhiza/web-monitor/web-api-read-only.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
