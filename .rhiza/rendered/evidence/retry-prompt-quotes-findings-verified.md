# retry-prompt-quotes-findings-verified

evidence · `01M7KXGGB5BFWBWXTXB7KF6GY7` · actor [claude-code-2bfb9046](../actor/claude-code-2bfb9046.md)

## result

PASS

## statement

At d24142b31 re-ran the Checks: block: go test -trimpath ./internal/api -run TestLeafPrompt (exit 0; TestLeafPromptQuotesFindings passes: a multi-line reason and one over 200 chars with leading/trailing spaces quoted byte for byte under 'Finding <id> (...)' heads, the inspect-the-working-tree and address-every-finding-exactly-as-written lines, the run status pointer kept; with no findings the failure detail is quoted and no findings list; a first attempt's prompt is exactly the rendered role prompt), go test -trimpath ./internal/run -run TestLeaves (exit 0; TestLeavesStepRetries finds the second builder's spec carrying both journalled findings verbatim and the third's spec the verifier failure with no findings), rhiza ship generate (exit 0, 0 changed) and rhiza render (exit 0, only this verifier's own actor rendered). rhiza show build-retries has the updated statement and acceptance naming the verbatim findings. Diff touches only the listed files plus rendered/actor records; findingLines, the checks step's fixer prompt and verifier journalling are unchanged; no TODOs.

## Links

### Out

- verifies [retry-prompt-quotes-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-prompt-quotes-findings.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
