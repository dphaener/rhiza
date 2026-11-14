# web-roadmap-milestone-data-qa-3

evidence · `01M7BY49CQN1S2PG3CKZ38NGMK` · actor [claude-code-bc270773](../actor/claude-code-bc270773.md)

## result

PASS

## statement

QA at 08472f9c4, fresh scratch build on a clone of the branch in a two-project root: rhiza milestones --json gave exit 0, 36 milestones, 6 finished, every change set with a non-empty intent (the new field), and rhiza api call milestones printed the same bytes. Each milestone's state word and each done point's met flag matched the rules (checked with jq). In the second project: a milestone with only a withdrawn change set was designing, one with a landed and a withdrawn one was building, and the finished count was 0. Concurrent /api/milestones loads all returned 200 and the daemon did not crash. Rough edge: a withdrawn change set's title is taken from its intent, not its --title field.

## Links

### Out

- verifies [web-roadmap-milestone-data](../capability/rhiza/web-monitor/web-roadmap-milestone-data.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
