# plan-permission-block-inbox

source · `01M6ZQF1C9N5S1JF26YQ1XR8MJ`

## quotes

```json
[]
```

## raw_sha256

4924598435a9fbcf7d07479acb9d7a6fd80d295aefafbca12ffc6261420ec66e

## session

779f56da-f75a-421d-a31c-01a1bf72ebf8

## summary

Autonomous plan run, no owner. The planner traced the bug: Claude Code's Notification permission_prompt hook sets a session blocked/permission and no hook fires when the prompt is answered, so session.Resolve's strict hook precedence hides a newer Rex idle and session message waits until its 10m timeout. inboxItems drops every block in autonomous runs and lists only question blocks otherwise, and Stalled never counts blocked sessions, so a permission prompt never reaches the inbox. The planner raised four questions (intent served, inbox item, how a stale block clears, what session message prints while waiting). The arbiter answered with Decision permission-block-inbox-and-clear: serve milestone-30; a permission-blocked session of a live run is an inbox item of kind permission at once, autonomous runs too; a newer non-blocked Rex report clears a permission block (question and auth blocks keep strict hook precedence); session message prints to stderr why it waits and its timeout error names the last status. Three leaves follow: permission-block-clears, inbox-permission-items, message-says-why-waiting.

## title

Plan: permission-blocked runs in the inbox, and session message on them (autonomous)

## Links

### In

- derived-from [permission-block-inbox-and-clear](../decision/permission-block-inbox-and-clear.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
