# health-across-projects-verified

evidence · `01M6SSAFPVDJHKBRBTVCFTDXGK` · actor [claude-code-3121443a](../actor/claude-code-3121443a.md)

## result

PASS

## statement

Verifier at 67772f59: read the diff (internal/api/health.go adds HealthRequest all/project, a project tag on HealthPart, machine parts once and project parts per scope via inboxScopes; describe-golden updated). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|describe-golden|inbox|roadmap)': ok, with health-projects.txtar PASS. That script builds a root rhiza.toml with projects alpha and beta. At the root, rhiza health --json shows daemon, disk, usage and machine once each, inbox and check twice, and 5 parts tagged with each project. A junk log file in beta makes rhiza health exit 1 with bad check beta and alpha not bad. rhiza health --project alpha --json and rhiza health --json run inside beta both give one inbox part and no project tag. --all inside beta covers both projects, and --project nope fails with not registered. go test -trimpath -count=1 ./internal/api -run 'Web|Health|Deploy' passed; the web caller passes a nil request and keeps the one-project report. go vet and gofmt are clean. The unscoped internal/api package run hit its 10m timeout while three other api.test processes loaded the 4-core machine; CI runs the full suite.

## Links

### Out

- verifies [health-across-projects](../capability/rhiza/orchestration/health-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
