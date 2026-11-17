# retry-prompt-quotes-findings

capability · `01M7KWHED131Y2GBPMA3KT4XK4` · parent [build-workflow](../build-workflow.md)

## acceptance

go test -trimpath ./internal/api -run TestLeafPrompt and go test -trimpath ./internal/run -run TestLeaves pass with tests that find a retry builder's prompt holding the previous attempt's two finding reasons byte for byte (one multi-line, one over 200 characters, neither trimmed nor joined), the inspect-the-working-tree line and the address-every-finding-as-written line; (failure case) an attempt that failed with no findings journaled (verifier ended with neither) gives a prompt quoting that attempt's failure detail and no empty findings list, and a first-attempt builder's prompt carries no findings section; rhiza show build-retries names the verbatim findings.

## files

internal/run/leaves.go, internal/run/leaves_test.go, internal/api/run_leaves.go, internal/api/run_leaves_test.go, build-retries (model node, via rhiza set)

## must_not

Change the fixer's findingLines shape or the checks step's fixer prompt; drop the rhiza run status --json pointer; change how a verifier journals findings

## statement

A leaves step's retry builder (attempt past 1 within its round) gets a prompt that quotes, verbatim and in full, each finding the previous attempt's verifier journaled on the leaf (or that attempt's failure detail when it journaled none), and tells it to inspect the working tree and the earlier attempts' commits before writing and to address every finding exactly as written, in one fix commit with the leaf's trailer; the landed build-retries node's statement and acceptance are updated to say so.

## Links

### Out

- derived-from [verbatim-retry-plan](../../../../../source/verbatim-retry-plan.md)

### In

- verifies [retry-prompt-quotes-findings-verified](../../../../../evidence/retry-prompt-quotes-findings-verified.md)
- verifies [retry-prompt-quotes-findings-qa-1](../../../../../evidence/retry-prompt-quotes-findings-qa-1.md)
- verifies [retry-prompt-quotes-findings-qa-2](../../../../../evidence/retry-prompt-quotes-findings-qa-2.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
