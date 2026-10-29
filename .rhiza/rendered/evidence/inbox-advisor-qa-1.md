# inbox-advisor-qa-1

evidence · `01M62KEDRDYN6RJE05MK12BC1N` · actor [claude-code-55549e75](../actor/claude-code-55549e75.md)

## result

PASS

## statement

QA with a build of the branch head (fake claude via the real hook path, fake rex): a normal run's agent blocked on 'Which colour?'; an advisor was launched with the question and 'follow rhiza guide advisor' in its prompt, journalled an advisor and a recommended entry, and rhiza inbox showed 'recommendation: Use blue' and 'reasons: Blue is the brand colour' beside the question while the run stayed blocked; nothing was typed at the asker and no Decision was written. With the advisor exiting 5 the question showed no recommendation (text or --json) and the run stayed blocked. rhiza guide advisor prints the recommend-only guide.

## Links

### Out

- verifies [inbox-advisor](../capability/rhiza/orchestration/inbox/inbox-advisor.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
