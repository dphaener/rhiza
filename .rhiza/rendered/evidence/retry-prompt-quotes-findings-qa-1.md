# retry-prompt-quotes-findings-qa-1

evidence · `01M7KZEE8FG9KH3MBQVCJ050J3` · actor [claude-code-805f61bc](../actor/claude-code-805f61bc.md)

## result

FAIL

## statement

QA at cf68954a8, scratch-built rhiza driving a leaves run with fake rex/claude: verifier journalling two findings (multi-line, 250+ chars with leading/trailing spaces) then exiting 0 gives the attempt-2 builder a prompt quoting both byte for byte, the inspect-the-working-tree and address-every-finding-as-written lines; a no-findings failure quotes 'verifier failed: exit 1'; attempt 1 has no retry section; three failures give no fourth builder and the leaf fails. FAIL: a verifier that journals findings and then exits 1 gives a retry prompt saying it journalled no findings (qa finding f1, internal/run/leaves.go:478).

## Links

### Out

- verifies [retry-prompt-quotes-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-prompt-quotes-findings.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
