# conduct-skill-replan-2026-11-13

source · `01M7B5N556EWZACW422S6RA5NW`

## quotes

```json
[]
```

## raw_sha256

68261fbfef41bad7ea5adba993aac74eb1f84adca97022c10bb16c89f0b0e931

## session

fb110e8b-6224-485a-a980-01a1d6583283

## summary

Autonomous re-plan run (01M7B5FZPG1F76KVT8866QBJMM) of change set 01M79KGZ (A rhiza-conduct skill conducts), milestone 36, already planned and accepted with five leaves (conduct-skill, check-plan-code-topics, review-three-checks, critic-on-major, conduct-plans-leaves) and Decision conductor-checks-gate-pick. There was no owner and no question was needed. The research pass re-read the slice and the code: every test, function and file the acceptances name still exists (TestGuideInstall, TestGuideInstallRoleSkills, TestReviewIsAnIndex, TestBuildReviewCommandBudget, TestDaemonPicks in daemon_picks_test.go, daemonCandidates' hasLeaves skip in internal/api/daemon.go, run_review.go's checks-step-only --check validation, workflow.go's Checks on checks steps); the sibling change sets run-spawn, guide-contract (accepted) and the conductor default (proposed) are where the plan expects them; rhiza overlaps is empty. One gap was fixed: conduct-skill tells the conductor to spawn the check-plan and check-code roles, which check-plan-code-topics creates, and both leaves edit internal/guide/guide.go and internal/api/guide_install.go, so conduct-skill now depends on check-plan-code-topics. The leaves, acceptances and Decisions are otherwise unchanged.

## title

Re-planning the conduct skill change set, autonomous

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
