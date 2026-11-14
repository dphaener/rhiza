# skill-check-builder-qa-2

evidence · `01M7E3FQQN61K90C0D9B4HK2RH` · actor [claude-code-a8f5980f](../actor/claude-code-a8f5980f.md)

## result

PASS

## statement

QA cycle 2 at b0e93e8f3: fresh go build -trimpath outside PATH, fake claude and fake rex in an isolated HOME and TMPDIR. rhiza skill check '/rhiza-builder <mode> {leaf}' --permission-mode manual: good 3 of 3, exit 0; bare fails trailers by name, exit 1; quiet fails session-done, exit 1; fail (exits 1) 2 of 3, exit 1. Each launched the builder once and no verifier (earlier qa f1 fixed). --json equals rhiza api call skill.check; uninstalled /rhiza-builder refused naming rhiza guide install --agent claude; empty invoke, bypassPermissions and TMPDIR inside a git repo refused; leaked RHIZA_CHANGE/RHIZA_ACTOR/RHIZA_SKILL_CHECK_ONCE ignored; user repo git status and .rhiza/log unchanged, scratch removed. Rough edge: a pass prints run: failed ... skill not installed (minor finding).

## Links

### Out

- verifies [skill-check-builder](../capability/rhiza/workflows/workflow-authoring/skill-check-builder.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
