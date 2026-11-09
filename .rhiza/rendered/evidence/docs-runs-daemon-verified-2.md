# docs-runs-daemon-verified-2

evidence · `01M6ZJR7NZB046R64KSCQD4HY5` · actor [claude-code-5ffa6feb](../actor/claude-code-5ffa6feb.md)

## result

PASS

## statement

Verified docs/runs-and-daemon.md at 3ae15b55 (review fix, cycle 2) as verifier claude-code-5ffa6feb. The confirmed findings are fixed. (1) Autonomy defaults: the page now says only plan is owner and accept, signoff and ship are auto; that matches internal/api/daemon_settings.go, and a scratch-root dry-run printed 'Autonomy plan: owner from rhiza.autonomy.plan' and 'accept: auto'. (2) Per-workflow autonomy: a new section covers top-level autonomous = true|false, the [autonomy] table, declaring only one of them, and workflow-over-setting precedence. In a scratch root with [autonomy] signoff/ship = owner in a copied workflow, `rhiza run start --workflow mine --dry-run` printed 'ship: owner from workflow mine' and 'signoff: owner from workflow mine', and the undeclared stages came from the settings. Setting both forms made `workflow check` report 'autonomous and [autonomy] are both set'. (3) Root prerequisite: in a fresh repo `rhiza request` failed with 'is not inside a Rhiza root', and it worked after `rhiza root init <dir>` and `rhiza project add <dir>/myapp`, as documented. (4) `rhiza run request-changes <findings.json> <run>` matches `request-changes <file> [<run>]`. docs-links-resolve: `go test -trimpath -count=1 ./internal/doclinks/...` ok. docs-commands-exist (checked by hand with a scratch build): all 30 commands the page names answer --help, and the flags it cites exist (--autonomous, --foreground, --workflow, --kind, --ref, --force, --note). The prose reads plain, consistent with doc-prose-humanized. Left open as minor, unconfirmed findings: no link to docs/custom-workflows.md, and the dry-run snippet needs a current change set. The owner's read is still to come at review.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
