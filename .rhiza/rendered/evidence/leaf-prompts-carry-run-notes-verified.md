# leaf-prompts-carry-run-notes-verified

evidence · `01M7KXWT1YDQFS74PY04RW9WPB` · actor [claude-code-f954a143](../actor/claude-code-f954a143.md)

## result

PASS

## statement

At 7eabdc606 re-ran the Checks: go test -trimpath ./internal/api -run 'TestLeafPrompt|TestRunSpawn' (TestLeafPromptNotes, TestRunSpawn/Notes, TestRunSpawnCheck/CheckRoleNoNotes pass), go test -trimpath ./internal/run -run TestRenderSpawnBrief and -run 'TestBrief|TestSession|TestHandback', rhiza ship generate and rhiza render (scratch build): all exit 0; go vet and gofmt clean. Diff read: withRunNotes appends run.NotesSection (oldest first, '- <note>' lines under 'Notes from earlier sessions of this run:') to builder and verifier leaf prompts only, read fresh from the journal; spawnNotes passes notes into RenderSpawnBrief only for builder/verifier children, after the caller's brief; no notes gives no heading; contract section and run note command unchanged; only the leaf's listed files changed.

## Links

### Out

- verifies [leaf-prompts-carry-run-notes](../capability/rhiza/workflows/built-in-workflows/build-workflow/leaf-prompts-carry-run-notes.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
