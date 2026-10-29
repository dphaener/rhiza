# inbox-open-qa-1

evidence · `01M62KECVE7CJXA1XHHQHKQSQT` · actor [claude-code-55549e75](../actor/claude-code-55549e75.md)

## result

PASS

## statement

QA with a build of the branch head against the real Rex: a throwaway tab was opened in a Rex session and fixture journals pointed a question (tab_session/tab_block) and a sign-off run (process rex_session/rex_block) at it. rhiza inbox open on each, run with no REX_SESSION, printed 'Focused <block>' / opened:true and rex session inspect showed that tab's window become active. With the tab closed, both printed '<kind> run <id>: tab gone (Rex no longer has <block>)' with exit 0; an unknown item exited 2.

## Links

### Out

- verifies [inbox-open](../capability/rhiza/orchestration/inbox/inbox-open.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
