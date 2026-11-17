# conduct-docs-reply-header-drafted

decision · `01M7MJYPS37KXQPBQA58TVNR6F`

## quote

yes, pass those ideas to the planner

## quote_where

the owner to the orchestrator, 2026-11-17, handing the idea to this planner (planner-notes 01M7EAYT705); the owner's words cover only the hand-off, not this choice

## reasoning

internal/api/session_message.go's messageContent already heads a message with the sender's handle and ID, the target session, change set, step and time, but not the sender's session or a reply command; adding them changes session.message's output and its tests, which this change set's intent (the guide is enough on its own) does not cover. The spawn-hang plan ruled the same way on the orchestrator's other rex-teams ideas (spawn-fix-stays-blocking).

## statement

This change set stays guide-only: the orchestrator's rex-teams idea that rhiza session message prefix each message with the sending session and the exact reply command is a change to session.message's behaviour, outside this intent, so it is drafted as its own change set 01M7MJYAZTA54VGAG8F6NAA67Y (discovered-from this one), which also documents it in the conduct guide; nothing here builds it.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [conduct-guide-selfcontained-plan-2026-11-17](../source/conduct-guide-selfcontained-plan-2026-11-17.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
