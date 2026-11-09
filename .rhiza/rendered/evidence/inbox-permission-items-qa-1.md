# inbox-permission-items-qa-1

evidence · `01M701V792PTVRA5C46AM1XA19` · actor [claude-code-1b20f401](../actor/claude-code-1b20f401.md)

## result

PASS

## statement

QA cycle 1, built binary at 73540a2d in a scratch root with fake rex/claude shims: rhiza run start --foreground (owner-run) whose agent hook reported Notification permission_prompt 'Allow Bash: echo secret@example.com?' showed at once in rhiza inbox / --json as kind permission with run, step build, session, detail, tab, tab_block, tab_session and next 'rhiza inbox open <id>'; rhiza inbox open <id> focused the session's block (rex focus -s session:s1 block:b1); rhiza daemon --once --notify-command ran the hook once with kind permission JSON and RHIZA_INBOX_TITLE '...: a session waits on a permission prompt: Allow Bash: ...'; a newer rex idle report in the session record cleared it (Nothing waits on you); an --autonomous run listed the permission item too. Scoped tests inboxapp, api TestWebInbox, scripts inbox-permission/daemon-notify-hook/inbox-stalled/inbox-open pass.

## Links

### Out

- verifies [inbox-permission-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-permission-items.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
