# inbox-milestone-interview-qa-1

evidence · `01M6N3J3FWG17B3TPWW80RNDQM` · actor [claude-code-23c10d47](../actor/claude-code-23c10d47.md)

## result

PASS

## statement

QA at dea4e5b6, built binary, fake rex/claude: while a non-autonomous milestone plan run's interview was live, rhiza inbox listed one 'interview' item naming the run, change set and title, tab 'repo: Backlog empty · run <id> · interview', block/session of the tab and 'the milestone interview waits on you'; rhiza inbox open <id> ran rex block inspect then rex focus on that tab's block; after the run ended the inbox was empty; feature and autonomous milestone runs gave no item. rhiza daemon --once with the interview live sent one blocked Rex report 'Backlog empty: a milestone interview waits on you'; a second --once sent none.

## Links

### Out

- verifies [inbox-milestone-interview](../capability/rhiza/orchestration/inbox/inbox-milestone-interview.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
