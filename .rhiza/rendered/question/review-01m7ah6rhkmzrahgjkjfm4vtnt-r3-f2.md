# review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f2

question · `01M7AMAMVDCNQMNT4J45Q5QWVS`

## statement

Plan review round 3 left finding f2 (contradiction on arbiter-advisor-brief): Decision brief-files-and-invoke-placement gives every session its own subdirectory at <step>-<visit>/<session>, so that later visits or cycles never share one. This leaf names the arbiter's directory <arbiter step>/arbiter-<question id>, with no visit, and the advisor's advisor/advisor-<asking session id>. A session that asks the owner two questions starts two advisors (StartAdvisor in internal/api/run.go runs once per question asked), and both get the same directory, so the second overwrites the first's brief.md and slice.json. Key the advisor's directory on something unique to each advisor, such as the question or the advisor's own session, and keep the arbiter in line with the <step>-<visit> layout, or record why these two sessions differ. The critic confirmed it: Real: the decision gives every session its own <step>-<visit>/<session> directory so that no two sessions share one. arbiter-advisor-brief names the advisor's directory advisor/advisor-<asking session id>. In internal/run/engine.go, block() journals a new block whenever the session is not already blocked, so a session can block again after an answer. advise() then calls StartAdvisor with the same asker session each time (engine.go:1090), and both advisors get the same directory, the second overwriting the first's brief.md and slice.json. The arbiter's <arbiter step>/arbiter-<question id> also drops the decision's -<visit>. The leaf contradicts the decision: major, as the contradiction kind is graded. How should it be settled?

## Links

### Out

- blocks [arbiter-advisor-brief](../capability/rhiza/workflows/workflow-runs/arbiter-advisor-brief.md)

### In

- answers [arbiter-advisor-dirs-per-session](../decision/arbiter-advisor-dirs-per-session.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
