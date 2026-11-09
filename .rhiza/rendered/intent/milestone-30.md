# milestone-30

intent · `01M6W5BP50YK3XRZTGRHBEC7EK`

## milestone

```json
30
```

## statement

Rhiza is ready to share: a developer with the link installs it, tries it on their own repo from the README and docs, and the repo is public under MIT.

## Links

### Out

- decided-by [m30-audience](../decision/m30-audience.md)
- decided-by [m30-plan-all-now](../decision/m30-plan-all-now.md)

### In

- serves change set `01M6ZQ588DTHT2T3MHK7STGR3Q` (Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.)

## History

- `01M6W5A2R4BXDGQGHNWGZNB3XY` Milestone 30: Rhiza is ready to share: install, README and docs, generated reference, public under MIT
