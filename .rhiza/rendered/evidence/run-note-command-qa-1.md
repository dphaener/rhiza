# run-note-command-qa-1

evidence · `01M7KZEE4H7P0WV6J5K2WQ2M3W` · actor [claude-code-805f61bc](../actor/claude-code-805f61bc.md)

## result

PASS

## statement

QA at cf68954a8, scratch-built rhiza (go build -trimpath): rhiza run note journalled with session/role/leaf from env and via --run prefix and api call run.note; run status human '## Notes' and --json notes oldest first with role and leaf; empty/space-only, \n, \r, 201-rune (200 é accepted), missing text and no RHIZA_RUN all exit 2 with nothing journalled; unknown run exit 4; --dry-run writes nothing. Scoped tests: internal/run -run TestNote, internal/api -run TestRunNote, TestScripts/run-note and describe-golden pass.

## Links

### Out

- verifies [run-note-command](../capability/rhiza/workflows/workflow-runs/run-note-command.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
