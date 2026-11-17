# leaf-prompts-carry-run-notes

capability · `01M7KWHEHXWHCRXW0CQVV9CF2J` · parent [build-workflow](../build-workflow.md)

## acceptance

go test -trimpath ./internal/api -run 'TestLeafPrompt|TestRunSpawn' and go test -trimpath ./internal/run -run TestRenderSpawnBrief pass with tests that journal two notes and find both, oldest first, in a later leaf's builder prompt, its verifier prompt and a spawned builder child's brief.md after the caller's brief; (failure case) a run with no notes gives prompts and briefs with no notes heading, and a spawned check-plan child's brief carries no notes.

## files

internal/api/run_leaves.go, internal/api/run_leaves_test.go, internal/api/run_spawn.go, internal/api/run_spawn_test.go, internal/run/brief.go, internal/run/brief_test.go

## must_not

Carry notes into check, critic or fixer prompts; change the contract section of brief.md; change the run note command

## statement

A leaves step's builder and verifier prompts, and the brief.md rhiza run spawn writes for a builder or verifier child, end with the run's notes so far, oldest first, one `- <note>` line each under a `Notes from earlier sessions of this run:` heading; with no notes they carry no such section.

## Links

### Out

- depends-on [run-note-command](../../workflow-runs/run-note-command.md)
- decided-by [run-notes-journaled](../../../../../decision/run-notes-journaled.md)

### In

- verifies [leaf-prompts-carry-run-notes-verified](../../../../../evidence/leaf-prompts-carry-run-notes-verified.md)
- verifies [leaf-prompts-carry-run-notes-qa-1](../../../../../evidence/leaf-prompts-carry-run-notes-qa-1.md)
- verifies [leaf-prompts-carry-run-notes-qa-2](../../../../../evidence/leaf-prompts-carry-run-notes-qa-2.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
