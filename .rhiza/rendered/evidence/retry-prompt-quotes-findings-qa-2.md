# retry-prompt-quotes-findings-qa-2

evidence · `01M7M1V7QXMM59ZQGAW8SXAZ6K` · actor [claude-code-91a0ded1](../actor/claude-code-91a0ded1.md)

## result

PASS

## statement

QA cycle 2 at c404b4730, scratch-built rhiza driving a leaves run with fake rex/claude: a verifier that journals two findings (multi-line with trailing newline, 200+ chars with leading/trailing spaces) and then exits 1 gives the attempt-2 builder a prompt quoting both byte for byte with the inspect-the-working-tree and address-every-finding-exactly-as-written lines (the cycle-1 qa f1 is fixed); a verifier ending with neither gives a retry quoting 'verifier failed: verifier ended with neither Evidence nor findings' and no findings list; attempt 1 has no retry section; three failures give no fourth builder, the leaf fails and attempt 3 quotes attempt 2's findings. TestLeafPrompt and internal/run TestLeaves pass.

## Links

### Out

- verifies [retry-prompt-quotes-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-prompt-quotes-findings.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
