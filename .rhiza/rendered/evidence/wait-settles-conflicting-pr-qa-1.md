# wait-settles-conflicting-pr-qa-1

evidence · `01M6AJB17P2GGFZ6XEY9Z1HG77` · actor [claude-code-6255e6c3](../actor/claude-code-6255e6c3.md)

## result

PASS

## statement

QA by hand at af6e4b54 with a scratch build (go build ./cmd/rhiza), in scratch repos with a bare origin and a fake gh on PATH whose pr view reports mergeable CONFLICTING until origin's branch holds main. Rendered-only conflict (both sides appended to .rhiza/rendered/index.md): rhiza ship check --settled printed 'not settled: the pull request conflicted; rebased onto main and pushed 0f55571' and exited 1; origin/cs/ship-it then contained origin/main; .git/rhiza/ship/<cs>.json held pushed = the new origin head; a second --settled --json call reported settled, passed and exited 0. Code conflict (both sides appended to code.txt): it exited 75, printing 'stop: rebasing cs/ship-it onto origin/main conflicts in code.txt; the rebase was aborted; answer question rebase-conflict-ship-it, then resume the run'; --json listed conflicts [code.txt] and question rebase-conflict-ship-it; the only new commit was 'Raise a question on the rebase conflicts', HEAD~1 and origin/cs/ship-it stayed at the pre-check head, the working tree was clean, and a second call exited 75 again without raising a duplicate question.

## Links

### Out

- verifies [wait-settles-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-settles-conflicting-pr.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
