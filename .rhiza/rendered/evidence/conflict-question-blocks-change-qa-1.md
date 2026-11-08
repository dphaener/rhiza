# conflict-question-blocks-change-qa-1

evidence · `01M6D2CZQ337RCC97WK04CFTVK` · actor [claude-code-1650fb1f](../actor/claude-code-1650fb1f.md)

## result

PASS

## statement

QA at 817b7c33 with a binary built from the branch (go build ./cmd/rhiza), in scratch repos with a bare origin and no fakes: an accepted change set 'QA ship' on worktree branch cs/qa-ship adding feature.txt, while main gained its own feature.txt. rhiza ship rebase --change <cs> exited 1, reported the feature.txt conflict and 'Raised question rebase-conflict-qa-ship.'; the one log entry held both node.create (question) and link.create type blocks to the change set. rhiza show listed 'blocks change set ... QA ship'; rhiza inbox --json listed it with kind blocker, blocks 'this change set'. Running ship rebase again left the .rhiza/log byte count and branch HEAD unchanged. rhiza question answer rebase-conflict-qa-ship 'Keep main' --change <cs> wrote the answer (show: 'answered by answer-rebase-conflict-qa-ship') and the inbox emptied. Upgrade path: a question raised by the merge-base binary had no blocks link and no inbox item; the new binary's ship rebase added only the blocks link (one new entry, no duplicate question), the inbox then listed it as a blocker, and a further run wrote nothing.

## Links

### Out

- verifies conflict-question-blocks-change (retired)

## History

- `01M6D12NNWQZHEWFBX60JS04AG` A run's rebase-conflict question blocks its change set, so the inbox and question answer see it
