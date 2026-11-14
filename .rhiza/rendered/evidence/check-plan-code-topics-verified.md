# check-plan-code-topics-verified

evidence · `01M7BHT7AVCT3DTKYHZRAW5MGQ` · actor [claude-code-a34e286e](../actor/claude-code-a34e286e.md)

## result

PASS

## statement

At a64eb85cc (leaf commits 9d24c14ff, a64eb85cc) re-ran the Checks block: go test -trimpath ./internal/guide -run 'TestCheckPlan|TestCheckCode|TestReviewIsAnIndex|TestBuildReviewCommandBudget' exit 0 (TestCheckPlanText, TestCheckPlanNamesNoRepoCommand, TestCheckCodeText, TestCheckCodeNamesNoRepoCommand, TestReviewIsAnIndex, TestBuildReviewCommandBudget PASS); go test -trimpath ./internal/api -run TestGuideInstallRoleSkills exit 0; go test -trimpath ./internal/guide exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (no tracked changes). Also ran go test -trimpath ./internal/api -run TestGuideInstall and TestScripts/guide and guide-install: ok. Diff read: check-plan.md merges the compliance criteria (under-built, scope creep, unaccounted files) with the model rules (code contradicts a node, missing constraint), with the role sections and shared blocks, tested under TopicCap; check-code.md is check-quality.md renamed (only title, task verb and plan-check reference differ); check-compliance, check-model and check-quality texts unchanged; review.md names check-plan, check-code, check-qa for a conductor run and the older four for ship-change, review and docs-review; guide_install.go installs skills/rhiza-check-plan and skills/rhiza-check-code beside the old role skills, checked by TestGuideInstallRoleSkills. Other touched files (guide.txtar, guide-install.txtar, contract_test.go, guide_test.go) are count/list updates the new topics require. No stubs or TODOs added; fakes only.

## Links

### Out

- verifies [check-plan-code-topics](../capability/rhiza/guide/agent-guide/check-plan-code-topics.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
