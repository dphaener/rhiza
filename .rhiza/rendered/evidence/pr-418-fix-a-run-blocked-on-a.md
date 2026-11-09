# pr-418-fix-a-run-blocked-on-a

evidence · `01M702XD8086RNEC5TXNB9AMN7` · actor [dphaener](../actor/dphaener.md)

## head

e45689667029691dd5930547670fe801333ef96b

## pr

https://github.com/dphaener/rhiza/pull/418

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #418 (https://github.com/dphaener/rhiza/pull/418) passed CI at e456896

## status

verified

## Links

### Out

- verifies [permission-block-clears](../capability/rhiza/runtime/agent-status/permission-block-clears.md)
- verifies [inbox-permission-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-permission-items.md)
- verifies [message-says-why-waiting](../capability/rhiza/orchestration/step-messages/message-says-why-waiting.md)
- verifies [status-precedence](../capability/rhiza/runtime/agent-status/status-precedence.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
