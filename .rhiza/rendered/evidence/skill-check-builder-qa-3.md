# skill-check-builder-qa-3

evidence · `01M7E42S1HPSH03RKCRZFAZY0N` · actor [claude-code-2248efc5](../actor/claude-code-2248efc5.md)

## result

PASS

## statement

QA cycle 3 at c696d215f: fresh go build -trimpath outside PATH, fake claude and fake rex, isolated HOME and TMPDIR, session env cleared. rhiza skill check '/rhiza-builder <mode> {leaf}' --permission-mode manual: good 3 of 3 exit 0 with run line 'stopped by skill check after the builder' (no install advice, earlier qa f2 fixed); bare fails trailers, quiet and fail (exit 1) fail session-done, each exit 1; every case launched the builder once and no verifier. --role planner, empty invoke, uninstalled /rhiza-nosuch (names rhiza guide install --agent claude), bypassPermissions and an unknown agent refused. User repo git status and .rhiza/log unchanged; scratch removed. TestScripts skill-check passes.

## Links

### Out

- verifies [skill-check-builder](../capability/rhiza/workflows/workflow-authoring/skill-check-builder.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
