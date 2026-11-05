# question-6b593167-1

question · `01M6JDV92Y0SKYFFNNNMHWTRC5`

## session

6b593167-fb44-4178-8462-01a1a4da96c7

## statement

Change set 01M6JDM (fix: request plans serve no milestone) serves no intent yet. It fits milestone-25 (01M6HKVRE1B7: 'every change set planned from the owner's queued requests ... is built and shipped by the daemon'). Serving it also lets the daemon build this fix. Which intent should it serve? / How does a request's plan reach a milestone? Today `rhiza request` takes only text and kind, and the planner picks any landed intent, often one under no milestone (e.g. the 'areas' umbrella), so the daemon skips it. / What stops a plan from ending with a change set that serves no milestone again (the failure here was silent: planned, accepted, then skipped)? / Should the requests TUI app (q) let the owner name the milestone when filing a request?

## Links

### In

- answers [request-plans-serve-a-milestone](../decision/request-plans-serve-a-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
