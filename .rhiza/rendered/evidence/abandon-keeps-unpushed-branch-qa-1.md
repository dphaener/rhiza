# abandon-keeps-unpushed-branch-qa-1

evidence · `01M6MG0669EHSMFM01S4WE744C` · actor [claude-code-de64ffbd](../actor/claude-code-de64ffbd.md)

## result

PASS

## statement

QA at cbdf7a93 with a binary built from the branch (go build ./cmd/rhiza) in a scratch repo with a bare origin, RHIZA_ROOT/CHANGE/RUN unset. Unpushed case: run A (command build workflow) committed on cs/abandon-built, pushed with git push -u origin; run B committed again without pushing; run C failed its check; rhiza run abandon C printed 'remove-worktree done' and 'delete-branch kept: it has commits merged into neither its upstream nor HEAD', and git branch -vv showed cs/abandon-built [origin/cs/abandon-built: ahead 1] with both built commits. Control: same flow without run B (branch fully pushed) printed 'delete-branch done' and the branch was gone. go test -trimpath ./cmd/rhiza -run 'TestScripts/run-abandon-(unpushed|unmerged)$' passed.

## Links

### Out

- verifies [abandon-keeps-unpushed-branch](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-keeps-unpushed-branch.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
