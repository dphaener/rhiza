# skill-check-builder-qa-1

evidence · `01M7E24TMCZ9S7N7PPDJQHTAVM` · actor [claude-code-98071bd0](../actor/claude-code-98071bd0.md)

## result

PASS

## statement

QA: scratch build of the branch (go build -trimpath, outside PATH) with a fake claude and fake rex in an isolated HOME. rhiza skill check "/rhiza-builder good {leaf}" --permission-mode manual: 3 of 3 gates pass, exit 0; bare (no trailers): fail trailers by name, exit 1; quiet (no session done): fail session-done, exit 1; crashing builder: 0 of 3, exit 1; --json equals rhiza api call skill.check output; uninstalled /rhiza-nosuch refused naming rhiza guide install --agent claude; empty invoke, --role planner, --dry-run, missing --permission-mode, bypassPermissions and an unknown agent refused; TMPDIR inside a git repo refused; leaked RHIZA_ROOT/RHIZA_CHANGE/RHIZA_ACTOR ignored; user repo git status and .rhiza/log unchanged, scratch removed unless --keep. Defect: a failing builder is retried up to 3 attempts, and a fixture verifier session launches after each attempt; only attempt 1 is graded (finding recorded).

## Links

### Out

- verifies [skill-check-builder](../capability/rhiza/workflows/workflow-authoring/skill-check-builder.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
