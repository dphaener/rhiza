# question-779f56da-1

question · `01M6ZQAWWE7K6CFRXKS16XA7QB`

## session

779f56da-f75a-421d-a31c-01a1bf72ebf8

## statement

Change set 01M6ZQ5 (permission-blocked run missing from inbox; session message hangs): which intent should it serve? / Should a session of a live run that is blocked on a permission prompt (hook kind=permission) be listed in rhiza inbox as soon as it blocks, including in autonomous runs, where the inbox lists no questions today? / Answering a permission prompt fires no Claude Code hook, so the session's hook status stays 'blocked permission' and session message waits until its timeout. How should a stale permission block clear? / While session message waits on a blocked session, what should it say?

## Links

### In

- answers [permission-block-inbox-and-clear](../decision/permission-block-inbox-and-clear.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
