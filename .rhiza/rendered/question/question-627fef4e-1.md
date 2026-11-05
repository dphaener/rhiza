# question-627fef4e-1

question · `01M6PQXYWRF7FEA88EHHHBHDDG`

## session

627fef4e-be7d-42fa-908e-01a1ad7d72e6

## statement

Change set 'Fix An owner-run plan's gate' serves no intent yet. Which should it serve? (Recent Fix change sets serve milestone-25, 'The backlog is empty… bugs and polish'.) / When an owner-run (non-autonomous) plan's gate auto-accepts because rhiza.autonomy.accept is auto, which rule does the accept use? Rule 'autonomous' today is allowed only inside an unfinished autonomous run. / Before an owner-run plan's gate auto-accepts (accept auto), must the plan's review have passed, as the autonomous gate requires (gateReview: a review round, no confirmed finding without a question, no open review question)?

## Links

### In

- answers [owner-run-gate-auto-accept](../decision/owner-run-gate-auto-accept.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
