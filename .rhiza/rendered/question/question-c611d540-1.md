# question-c611d540-1

question · `01M6VYDPA4Y02J9HAVXJP0Q232`

## session

c611d540-dc1b-4255-bc52-01a1b7e438f5

## statement

Change set 01M6VY7 (Scheduling as policy) must serve a milestone today: the plan's milestone-check fails otherwise, and the daemon skips it as 'serves no milestone'. Milestones 26–29 have no open change sets left. Which milestone should 01M6VY7 serve? / Where should the repo's scheduling rule for change sets that serve no milestone live? / Which values should the scheduling rule take, and what is the default? / What should the plan workflow's milestone-check step do once backlog work can ship?

## Links

### In

- answers [scheduling-policy-shape](../decision/scheduling-policy-shape.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
