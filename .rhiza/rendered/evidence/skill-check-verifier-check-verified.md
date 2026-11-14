# skill-check-verifier-check-verified

evidence · `01M7DYPQ3XFGV0P4PWBYZ0S6PE` · actor [claude-code-6a05e316](../actor/claude-code-6a05e316.md)

## result

PASS

## statement

At f8e3bc8a6 re-ran the Checks: block: go test -trimpath -p 2 -count=1 ./cmd/rhiza -run 'TestScripts/^(skill-check-roles|skill-check|describe-golden)$' ok; go test ./internal/api -run TestSkillCheck ok (TestSkillCheckRoles and its Definitions, Unused skill refused, Findings gate, Not started, Report subtests pass); ./internal/api -run 'Describe|Help|Schema|Decode|Registry|Operation' ok; go vet ./internal/api clean; describe-golden -update, rhiza ship generate and rhiza render left no diff to the commit's files; make self-check passes in a clean detached worktree. skill-check-roles.txtar, with a fake claude and fake rex, shows: a verifier that writes Evidence and says done passes 3 of 3 gates (session-done, evidence, check); one writing a finding instead fails 'evidence' by name and exits non-zero (ExitCode 1 asserted in TestSkillCheckRoles/Report); a check recording a round with rhiza run findings --check passes 3 of 3; one recording none fails 'findings' by name; --role planner is refused with 'want builder, verifier or check'; git status and .rhiza/log of the current repo are unchanged after every case, and no builder session starts. The diff touches the three named files plus the regenerated describe-golden, docs/reference/cli.md, and the builder's actor record and renders; no TODOs or stubs.

## Links

### Out

- verifies [skill-check-verifier-check](../capability/rhiza/workflows/workflow-authoring/skill-check-verifier-check.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
