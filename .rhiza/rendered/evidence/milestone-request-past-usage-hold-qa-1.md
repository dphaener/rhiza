# milestone-request-past-usage-hold-qa-1

evidence · `01M6N3J8RDN7CGB83CRE0SQ1YJ` · actor [claude-code-23c10d47](../actor/claude-code-23c10d47.md)

## result

PASS

## statement

QA at dea4e5b6, built binary, real detached runs with fake rex/claude: with rhiza.daemon.usage.session 1000 and a Claude transcript of 2000 tokens, a queued feature and a queued milestone request: paused, daemon --once started neither; --min-free 1000T (disk low) started neither; unpaused under the usage-reserve hold (inbox lists usage-reserve) it started only the milestone request's plan run and request ls showed the feature still queued.

## Links

### Out

- verifies [milestone-request-past-usage-hold](../capability/rhiza/orchestration/usage-governor/milestone-request-past-usage-hold.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
