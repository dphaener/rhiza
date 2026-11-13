# plan-prompts-one-line

capability · `01M79BH6T4DQAX8RJF1YTQC0AS` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/workflow/... passes with tests that plan.toml's and plan-milestone.toml's interview, review, critic and revise prompts equal the four one-liners (no `session done` text left in them) and the prepush-fix prompt is unchanged; go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-|bugfix-workflow-autonomous|arbiter-agent-tabs|status)' passes with fake agents matching the new prompts; the describe-golden script is regenerated only if its output changed; make self-check passes.

## statement

The builtin plan.toml and plan-milestone.toml interview, review, critic and revise prompts are one line each that carries the TASK variables and names both routes: `/rhiza-planner Plan change set {change} ({title}): follow rhiza guide planner.`, `/rhiza-plan-reviewer Review the plan of change set {change}: follow rhiza guide plan-reviewer.`, `/rhiza-plan-critic Judge the latest review round of the plan of change set {change}: follow rhiza guide plan-critic.` and `/rhiza-plan-reviser Fix the confirmed findings of the latest review round of the plan of change set {change}: follow rhiza guide plan-reviser.`; the session-done and command rules now in those prompts live in the topics, and the prepush-fix prompt is unchanged.

## Links

### Out

- depends-on [plan-role-skills-installed](plan-role-skills-installed.md)
- decided-by [plan-role-topics-shape](../../../../../decision/plan-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../decision/afk-prompts-shipped-default.md)
- derived-from [plan-role-topics-interview](../../../../../source/plan-role-topics-interview.md)

### In

- verifies [plan-prompts-one-line-verified](../../../../../evidence/plan-prompts-one-line-verified.md)
- verifies [plan-prompts-one-line-qa-1](../../../../../evidence/plan-prompts-one-line-qa-1.md)
- verifies [pr-533-plan-role-topics](../../../../../evidence/pr-533-plan-role-topics.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
