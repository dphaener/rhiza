# build-role-skills-installed-verified-5

evidence · `01M79CGQTFJHSWJRH7JM499BQM` · actor [claude-code-b18d95ad](../actor/claude-code-b18d95ad.md)

## result

PASS

## statement

At 90cf3007a, attempt 5, as verifier: go test -trimpath -count=1 ./internal/api -run GuideInstall passed, incl. TestGuideInstallRoleSkills (rhiza-builder, rhiza-verifier, rhiza-build-retry SKILL.md written with topic description and expanded body; second install reports them in Unchanged) and the pluginPaths written/unchanged round trip. go test -trimpath -count=1 ./internal/guide/... passed; TestBuildIsAnIndex asserts the build index names builder, verifier and build-retry each by skill (/rhiza-builder etc.) and `rhiza guide` route, and that the topics exist. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(guide|guide-install|build-evidence-commit)$' passed: guide.txtar checks guide ls lists build-retry, builder, verifier; guide-install.txtar checks the three new skill paths in the install output and their frontmatter/body against `rhiza guide <topic>`; build-evidence-commit.txtar's fake builder evals the actor-record git add line from `rhiza guide builder`. A scratch build's `rhiza guide build` prints the short index naming the three role topics and which role loads each.

## Links

### Out

- verifies [build-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-role-skills-installed.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
