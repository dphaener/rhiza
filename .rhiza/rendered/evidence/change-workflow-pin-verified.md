# change-workflow-pin-verified

evidence · `01M731BRPMBAAKZ3WYXTACBNFF` · actor [claude-code-9d186c3f](../actor/claude-code-9d186c3f.md)

## result

PASS

## statement

At f82e440f, go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-workflow-pin|describe-golden)$' passes. The change-workflow-pin testscript pins a draft change set with rhiza change workflow docs. Then change show --json shows "workflow": "docs" and run start --dry-run prints '- Workflow: docs'. --clear returns previous=docs, after which change show has no workflow and run start --dry-run falls back to the dispatch rules. The script refuses an unknown name ('no such workflow'), a name together with --clear, no name at all, and a landed change set ('landed'). Each stage passes fold-ok. go test -run TestChangeWorkflow ./internal/api passes; it covers the change.workflow operation's refusals (exit 2), dry run, pin, clear and the landed refusal (exit 3). internal/changeset, entry, store and render all pass, including the format-5 pin corpus. I read the diff: the pin entry is format 5, the fold applies it only in flight, and the operation refuses a live run through refuseLiveRun. The full internal/api run failed only on things this change does not touch: the driver poll idle CPU budget (5.7ms over 5ms on the shared host) and the existing .rhiza-root/usage.live.json leak check from the run-start tests.

## Links

### Out

- verifies [change-workflow-pin](../capability/rhiza/workflows/workflow-dispatch/change-workflow-pin.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
