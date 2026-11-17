# leaf-prompts-carry-run-notes-qa-1

evidence · `01M7KZEEDVECW1B6ZCNHDG8YGZ` · actor [claude-code-805f61bc](../actor/claude-code-805f61bc.md)

## result

PASS

## statement

QA at cf68954a8, same scratch leaves run: after a builder and a verifier ran rhiza run note, later builder and verifier prompts (same and next leaf, retries too) end with 'Notes from earlier sessions of this run:' and both '- <note>' lines oldest first; the first builder (no notes yet) has no heading; a spawned builder child's brief.md carries the notes after the caller's brief, a spawned check-plan child's carries none. Scoped tests internal/api -run 'TestLeafPrompt|TestRunSpawn', internal/run -run TestRenderSpawnBrief, TestScripts/run-spawn pass.

## Links

### Out

- verifies [leaf-prompts-carry-run-notes](../capability/rhiza/workflows/built-in-workflows/build-workflow/leaf-prompts-carry-run-notes.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
