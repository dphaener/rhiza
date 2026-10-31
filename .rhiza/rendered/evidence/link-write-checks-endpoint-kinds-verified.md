# link-write-checks-endpoint-kinds-verified

evidence · `01M686JHTN25TM9Z1AQ1WS32KB` · actor [claude-code-9479858b](../actor/claude-code-9479858b.md)

## result

PASS

## statement

Verified at 2fca47e. Read the diff: buildLink in internal/api/write_link.go refuses a link unless check.Allowed(type, fromEnd, toEnd), which is the same table internal/check/check.go uses. The refusal is a usage error, "link type %q does not go from %s to %s", raised before anything is written. go test ./cmd/rhiza -run 'TestScripts/(link|apply|check)$' passes. link.txtar runs store serves log with --json (code usage), plain, --dry-run and --record, refuses a change set serving a capability and accepts one serving an intent. apply.txtar refuses a whole apply with one bad link item and adds no entry. check.txtar and write_test.go cases that wrote bad links now use valid ones. go test ./internal/api -run TestLinkEndpointKinds passes: for every link type and every pair of end kinds, change sets included, buildLink agrees with check.Allowed and refuses with exit 2. By hand, with a scratch binary in a fresh repo and run env unset, rhiza link store serves log with --json, with --dry-run and with no flag, and rhiza link a serves b --record on landed capabilities, each exit 2 with 'link type "serves" does not go from capability to capability'. The log entry count stays at 12. In make check, fmt, vet and all Go tests pass. Its self-check fails only on stale .rhiza/rendered actor files and index.md, left by this run's actor records. The commit does not touch .rhiza.

## Links

### Out

- verifies [link-write-checks-endpoint-kinds](../capability/rhiza/model/typed-links/link-write-checks-endpoint-kinds.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
