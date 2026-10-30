# printout-gap-names-reason-verified

evidence · `01M65BZK796XHAWMF792ZS328N` · actor [claude-code-f9f7c666](../actor/claude-code-f9f7c666.md)

## result

PASS

## statement

Verified at 02ac599 in a detached checkout. Read the diff: RunSession gains no_tokens; transcriptTokens returns nil plus the reason ('still running' when the session record has not exited, else 'transcript unreadable: <warning>'); runInfo adds only sessions with tokens to the totals; RunInfo.Human, which both run status and the final printout of run start/sign-off/abandon use, prints 'no token use shown, <reason>' in place of 'no token use recorded'. Ran go test ./cmd/rhiza -run 'TestScripts/(printout-gap-names-reason|run-status-tokens|describe-golden)$' -count=1: ok. The new printout-gap-names-reason.txtar runs a three-session run, strips the journaled tokens, deletes one transcript and the exit report of the review session, and asserts run status --json gives no_tokens 'still running' for review and 'transcript unreadable: no transcript of session ...: looked in ...' for check, exactly one tokens object, and totals equal to the build session alone (10/20/30/40, helpers 1/2/3/4); the human run status names both reasons, never says 'no token use recorded', and its Totals line is the build session's only. describe-golden was regenerated because a response field (no_tokens) was added and a desc changed; the golden test passes. make check (fmt-check, vet, test, self-check) exits 0.

## Links

### Out

- verifies [printout-gap-names-reason](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-gap-names-reason.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
