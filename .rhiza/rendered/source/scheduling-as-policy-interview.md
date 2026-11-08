# scheduling-as-policy-interview

source · `01M6VYJBF4EMJPMNX947JYGKE1`

## quotes

```json
[]
```

## raw_sha256

b94f0d384cc185b70b3bdfc6542d1af8be180f0772e3424eaef5c212f9d0d26d

## session

c611d540-dc1b-4255-bc52-01a1b7e438f5

## summary

Autonomous plan run for change set 01M6VY7 (Scheduling as policy), which follows decision workflows-are-the-rules. No owner took part. The planner read the daemon's candidate selection (daemonCandidates in internal/api/daemon.go skips every accepted change set that serves no milestone as SkipNoMilestone), the plan workflow's milestone-check step (internal/workflow/builtin/plan.toml, internal/api/run_milestone.go), the request --milestone plumbing and the plan guide. It raised four questions in one call: which milestone this change set serves, where the scheduling rule lives, which values it takes, and what milestone-check does once backlog work can ship. The arbiter answered with Decision scheduling-policy-shape. (1) Serve milestone-29 directly, with no new milestone. (2) The rule lives in a new checked-in .rhiza/policy.toml under [schedule], so every daemon reads the same rule, and a missing file or key means the default. (3) backlog = after (default) | before | skip, oldest accepted first inside the backlog; an unknown value is a rhiza check problem, and the daemon then skips backlog work with the bad value in the skip detail. (4) milestone-check stays. It fails only for a run started with --milestone whose change set does not serve that milestone, or when backlog is skip. The planner then wrote four leaves: the policy file and its check, the daemon's backlog picks, the relaxed milestone-check, and the plan guide's milestone wording.

## title

Scheduling as policy: autonomous plan interview (arbiter answers)

## Links

### In

- derived-from [scheduling-policy-shape](../decision/scheduling-policy-shape.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
