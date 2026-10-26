# plan-task-1

source · `01M5SWRGMNEQC0N78KXY7RJSH4`

## quotes

```json
[]
```

## raw_sha256

619b89f2d41f231169c63419862b23bfee059ba2328ae3aeb8d91e00a2b43cea

## summary

The shipit work plan for task 1: turn the repo into a Go module with a thin genesis binary, a testscript harness against hermetic real git, a make check gate, one CI job and a short CLAUDE.md, delivered as one PR with green CI. Decided (D1-D15, agent wording accepted by the owner in the interview): the binary only prints a version, a short usage with no arguments, and fails on an unknown command, stdlib only, with task 6 replacing the dispatch; cmd/genesis stays a thin entry point and domain packages go under internal/; make check lints with a gofmt check and go vet, no golangci-lint; the module declares go 1.25.14; testscript runs in-process with a hermetic HOME and git config, scripts in cmd/genesis/testdata/script; all of .shipit/ is gitignored; work on a branch with a PR whose CI must be green, no git-spice; TDD with scenarios seen failing first; the version is stamped by ldflags from git describe, falling back to build info and then dev; Makefile targets including install built in place, never copied; a short CLAUDE.md; CI is one Ubuntu job running make check; out of scope are domain packages, CLI libraries, .genesis/ and genesis check in CI, and release workflows; a scenario proves hermetic git; acceptance is make check passing, an unformatted file failing it, a stamped version, and green CI. Nothing was revised after approval and nothing was left open.

## title

Shipit plan, milestone 1 task 1: repo foundation

## Links

### In

- derived-from [t1-d1-version-usage](../decision/t1-d1-version-usage.md)
- derived-from [t1-d2-cmd-internal-layout](../decision/t1-d2-cmd-internal-layout.md)
- derived-from [t1-d3-gofmt-vet-lint](../decision/t1-d3-gofmt-vet-lint.md)
- derived-from [t1-d4-go-module-version](../decision/t1-d4-go-module-version.md)
- derived-from [t1-d5-testscript-in-process](../decision/t1-d5-testscript-in-process.md)
- derived-from [t1-d6-shipit-gitignored](../decision/t1-d6-shipit-gitignored.md)
- derived-from [t1-d7-branch-pr-ci](../decision/t1-d7-branch-pr-ci.md)
- derived-from [t1-d8-tdd-testscript](../decision/t1-d8-tdd-testscript.md)
- derived-from [t1-d9-version-stamping](../decision/t1-d9-version-stamping.md)
- derived-from [t1-d10-makefile-targets](../decision/t1-d10-makefile-targets.md)
- derived-from [t1-d11-short-claude-md](../decision/t1-d11-short-claude-md.md)
- derived-from [t1-d12-ci-one-job](../decision/t1-d12-ci-one-job.md)
- derived-from [t1-d13-task-1-out-list](../decision/t1-d13-task-1-out-list.md)
- derived-from [t1-d14-hermetic-git-scenario](../decision/t1-d14-hermetic-git-scenario.md)
- derived-from [t1-d15-task-1-acceptance](../decision/t1-d15-task-1-acceptance.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
