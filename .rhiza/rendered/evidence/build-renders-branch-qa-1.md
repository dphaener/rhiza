# build-renders-branch-qa-1

evidence · `01M66DTZQGHJGM4RAY0S8W88BG` · actor [claude-code-fa0d8666](../actor/claude-code-fa0d8666.md)

## result

PASS

## statement

QA at 29929ba with this commit's built rhiza: go test TestScripts/build-workflow and run-ship-change(-fail) pass. In a scratch repo on a branch, committed a capability, then left a landed actor record (rhiza add actor --record) and an Evidence entry uncommitted, so rhiza check failed (rendered folder stale/missing). Ran build.toml's [step.render] command verbatim with RHIZA_CHANGE set: exit 0, one new commit 'Render after the build' with 'Change: <change id>', holding both log entries and the rendered actor pages and index; git status empty and rhiza check prints ok afterwards.

## Links

### Out

- verifies [build-renders-branch](../capability/rhiza/workflows/build-renders-branch.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
