# permission-block-inbox-and-clear

decision · `01M6ZQDYSJAT1RMJBA2DGAHWCA`

## reasoning

Fixes serve the open milestone: request plans link to the open milestone (01M6JDMJ), and this bug came up in milestone 30's pre-publish audit run. The inbox (internal/api/inbox.go inboxItems) drops every block in autonomous runs and lists only kind 'question' blocks otherwise. Its comment says only questions are the orchestrator's, and the arbiter decides them. A permission prompt needs a keypress in the terminal, so neither the arbiter nor the orchestrator can answer it. Leaving it out is how run 01M6ZMRG sat blocked 36+ min with 'Nothing waits on you'. Stalled() also never counts blocked sessions, so no other path brings it up. internal/session/status.go Resolve gives hooks strict precedence over Rex. Answering a permission prompt fires no Claude Code hook (claudecode.go maps only Notification permission_prompt to blocked), so the stale hook report hides a fresh Rex idle and typePointer (session_message.go) waits until its 10m timeout. That is the hang. The run already follows rex events (session_run.go recordRex), so a newer non-blocked Rex report is real evidence that the prompt closed. Limiting the override to permission blocks keeps precedence as it is for everything else and adds no guessing. session message printed nothing while it waited. Telling the caller what it waits on and how to clear it matches the change set's expectation that it deliver to an idle session or time out with an error.

## statement

(1) Change set 01M6ZQ5 serves milestone-30 (01M6W5BP50YK3XRZTGRHBEC7EK), the open milestone. (2) Yes: a live run's session whose resolved status is blocked with kind permission is an inbox item as soon as it blocks, in autonomous runs too, since no arbiter can answer a terminal prompt. Give it its own kind (permission). The item names the run, step, session, the prompt (the hook's detail message) and its tab, and says how to answer: rhiza inbox open <item> focuses the tab, then answer the prompt there. The notify hook and the web UI show it like the other kinds. Autonomous runs still list no questions. (3) A permission block clears when a newer report says the prompt is gone: any later hook report, as now, or a Rex report for the session that is newer than the hook's blocked-permission report and is not blocked. Then Resolve uses that Rex report (idle or working) in place of the stale hook report. With no newer report the block stands. Do not clear on a timer or by guessing. Question and auth blocks keep strict hook precedence. (4) While session message waits, it prints to stderr why it waits, once at the start and again whenever the resolved status changes. For a blocked session that line names the block's kind and detail and how to answer, e.g. 'session 411ae041 is blocked on a permission prompt (<detail>): answer it in its tab (rhiza inbox open <item>); waiting up to 10m for its prompt'. Working or starting sessions get 'session <id> is working; waiting up to <timeout> for its prompt'. The timeout conflict error names the last status. JSON output and the 10m default stay as they are.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-779f56da-1](../question/question-779f56da-1.md)
- derived-from [plan-permission-block-inbox](../source/plan-permission-block-inbox.md)

### In

- decided-by [permission-block-clears](../capability/rhiza/runtime/agent-status/permission-block-clears.md)
- decided-by [inbox-permission-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-permission-items.md)
- decided-by [message-says-why-waiting](../capability/rhiza/orchestration/step-messages/message-says-why-waiting.md)
- decided-by [status-precedence](../capability/rhiza/runtime/agent-status/status-precedence.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
