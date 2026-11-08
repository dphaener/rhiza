# bugfix-workflow-autonomous-verified-2

evidence · `01M6WXTCDWKRHKDFSBE12VDB6Z` · actor [claude-code-92834dfb](../actor/claude-code-92834dfb.md)

## result

PASS

## statement

At 4c9e4641 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/bugfix-workflow-autonomous$' -count=1: ok. Read the script and diff: fake rex, claude and gh (gh answers pr view from a bare origin); repo .rhiza/workflows/bugfix.toml and bugfix-ship.toml declare autonomous = true and dispatch.toml maps request/ready kind bug to them; git config rhiza.autonomy.plan owner and rhiza.autonomy.signoff owner. The daemon starts the bug request's plan run with --autonomous, the run succeeds with Autonomous: yes and its gate journals one accepted entry (rule autonomous); the daemon then starts the bugfix-ship run --autonomous, which signs off by rule autonomous (no actor), pushes, passes rhiza ship check through the fake gh, and runs rhiza ship land: the change set reads landed and the run branch head equals origin's cs/fix-crash. rhiza inbox --json holds no item after the daemon start, after the plan run, and after the ship run. This fixes the earlier finding that landing was never exercised.

## Links

### Out

- verifies [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
