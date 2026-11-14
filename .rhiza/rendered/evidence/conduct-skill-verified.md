# conduct-skill-verified

evidence · `01M7BJ2HRBG03QG9D74EQX8HMT` · actor [claude-code-6eeb8054](../actor/claude-code-6eeb8054.md)

## result

PASS

## statement

At 93f17e390 re-ran the Checks block: go test -trimpath ./internal/guide -run TestConduct, go test -trimpath ./internal/api -run TestGuideInstall, go test -trimpath ./cmd/rhiza -run TestScripts/guide, rhiza ship generate (0 changed) and rhiza render, all exit 0; full ./internal/guide also passes. Read the diff: internal/guide/text/conduct.md has Task, Expected outcome, Read first, Must do, Must not do, Context in order; Must do names rhiza apply for missing leaves (arbiter rule), rhiza run spawn with builder, verifier, check-plan, check-code, check-qa, critic (only on a blocker or major round) and fixer, rhiza run findings --check, rhiza run review-check, rhiza land, rhiza render, rhiza session done; it carries the shared command-budget, never-list, session-done and contract-ref blocks and never-do items merge or push, edit .rhiza/log, pass your own gate. TestConduct asserts sections, phase order, shared blocks, never-list and the topic token cap; TestGuideInstallRoleSkills asserts skills/rhiza-conduct/SKILL.md written and unchanged on reinstall. rhiza guide conduct from a scratch build renders the topic, and rhiza run spawn --help confirms --role/--leaf/--brief. Extra files (guide_test.go list order, guide.txtar and guide-install.txtar counts, .rhiza actor record and render) are required follow-ons of the new topic. No TODOs or stubs.

## Links

### Out

- verifies [conduct-skill](../capability/rhiza/guide/agent-guide/conduct-skill.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
