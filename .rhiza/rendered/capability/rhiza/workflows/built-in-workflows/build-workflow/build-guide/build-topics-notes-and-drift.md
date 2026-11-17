# build-topics-notes-and-drift

capability · `01M7KWHENJAYXZ8T9MQKCWP4MP` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with tests that find rhiza run note in the builder and verifier topics, the update-landed-nodes-or-raise-a-question rule in the builder topic, the verbatim-findings line in build-retry, and the verbatim-findings retry brief in conduct; (failure case) a negative test asserts the verifier topic still carries no edit-code instruction and rhiza run note appears in no check or critic topic; every topic stays under the guide size cap; go test -trimpath ./internal/api -run TestGuideInstall passes.

## files

internal/guide/text/builder.md, internal/guide/text/verifier.md, internal/guide/text/build-retry.md, internal/guide/text/conduct.md, internal/guide/builder_test.go, internal/guide/verifier_test.go, internal/guide/build_retry_test.go, internal/guide/conduct_test.go, builder-topic, verifier-topic, build-retry-topic and conduct-skill (model nodes, via rhiza set)

## must_not

Name a Go, make or repo-specific command in the builder text; change other topics; change any code outside internal/guide

## statement

The builder and verifier topics tell their session to end with one `rhiza run note` line of a learning a later leaf can use (when it has one) before session done; the builder topic tells a builder whose leaf changes what a landed node says to update that node with rhiza set in the change set, or raise a question when it cannot; the build-retry topic says the findings are quoted verbatim in the prompt; the conduct topic says a builder retry's brief quotes the verifier's findings verbatim and tells it to inspect the working tree; the landed builder-topic, verifier-topic, build-retry-topic and conduct-skill nodes are updated to match.

## Links

### Out

- depends-on [run-note-command](../../../workflow-runs/run-note-command.md)
- decided-by [run-notes-journaled](../../../../../../decision/run-notes-journaled.md)

### In

- verifies [build-topics-notes-and-drift-verified](../../../../../../evidence/build-topics-notes-and-drift-verified.md)
- verifies [build-topics-notes-and-drift-qa-1](../../../../../../evidence/build-topics-notes-and-drift-qa-1.md)
- verifies [build-topics-notes-and-drift-qa-2](../../../../../../evidence/build-topics-notes-and-drift-qa-2.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
