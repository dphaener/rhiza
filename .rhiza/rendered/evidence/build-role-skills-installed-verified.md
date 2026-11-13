# build-role-skills-installed-verified

evidence · `01M79AA187V8FV9SWC1BSB3KDY` · actor [claude-code-b7e3213f](../actor/claude-code-b7e3213f.md)

## result

PASS

## statement

Verified at 7ce66901a in a detached checkout. Re-ran the commit's Checks: block and the acceptance, all exit 0: go test -trimpath ./internal/api -run 'GuideInstall|SkillFile|Plugin|Status' (TestGuideInstallRoleSkills checks that rhiza-builder, rhiza-verifier and rhiza-build-retry SKILL.md are the frontmatter with the topic's description plus the expanded body, and that a second install reports them unchanged); go test -trimpath ./internal/guide/... (TestBuildIsAnIndex checks that build names builder, verifier and build-retry by their /skill and `rhiza guide` routes and stays under 250 tokens); go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install|build-evidence-commit|status|docs-commands-exist)$' (guide ls has the three topics, install output lists the new skill paths with 10 written, and build-evidence-commit's fake builder greps `rhiza guide builder`); rhiza ship generate, rhiza render (tree clean afterwards) and make self-check. Read the diff: guide_install.go writes the three skills beside the others, build.md is a short index saying which role loads which topic, the only log entries staged are the builder's own actor record, and no stubs or TODOs.

## Links

### Out

- verifies [build-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-role-skills-installed.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
