# milestone-wait-owner-steps-only

decision · `01M7677NAX95EK183R1B84B5EY`

## reasoning

Intent: the two change sets that built and fixed the milestone request queue (01M6MYFH4Q 'A milestone request automatically launches an interactive Claude planning session', 01M6XV4WGS 'Request plan runs start outside maxRuns') both serve milestone-25, whose statement is that the daemon builds and ships the owner's queued requests on its own so the owner can start on new things; a request stuck for an hour behind an unattended review/CI tail is exactly that gap. Rule: the change set's intent states the one-at-a-time rule exists so the owner holds one interview at a time, and asks to count only runs at an owner-interactive step. The plan-milestone workflow (rhiza workflow describe plan-milestone) has one owner interview as its start step, then agent and command steps (milestone-check, review, critic, check, revise, gate, open-pr, wait-pr, pr-green, rebase, merge, close-out) that need no owner unless a step is paused on the owner. Keeping the unjournaled-start case blocking (as both functions already do) avoids two interviews opening in the same tick. Leaving maxRuns as the cap keeps the fix from 01M6XV4WGS intact.

## statement

Change set 01M767125EJM98D7H4ME2T67H5 serves @milestone-25. A queued milestone request waits only while another milestone request's plan run holds the owner: its start is still unjournaled (within daemonStartGrace, so it is about to open its interview), or it is live and its current step is owner-interactive: the workflow's interview (start) step, or a step paused waiting on the owner (a gate or sign-off under owner autonomy). A live milestone plan run on any other step (milestone-check, review, critic, check, revise, open-pr, wait-pr, CI, rebase, merge, close-out) blocks nothing, so the next milestone request starts as soon as the previous one's interview has passed. daemonMilestoneWaits and daemonPlanRequests' milestoneLive apply this same rule, and rhiza.daemon.maxRuns still caps every start.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a77d1ef4-1](../question/question-a77d1ef4-1.md)
- derived-from [plan-fix-queued-milestone-request](../source/plan-fix-queued-milestone-request.md)

### In

- decided-by [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- decided-by [milestone-wait-owner-steps](../capability/rhiza/orchestration/daemon/plan-requests/milestone-wait-owner-steps.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
