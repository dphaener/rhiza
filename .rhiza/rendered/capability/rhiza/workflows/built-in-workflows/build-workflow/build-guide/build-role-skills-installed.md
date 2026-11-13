# build-role-skills-installed

capability · `01M7974G8PCWZ6YB020VYE5FRB` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/api -run GuideInstall passes with the three new skill files written with their topic's description and expanded body, and re-running install reports them unchanged; go test -trimpath ./internal/guide/... asserts the build index names builder, verifier and build-retry and both their skill and `rhiza guide` routes; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install|build-evidence-commit)' passes with the new topics in guide ls, the new skill paths in the install output, and build-evidence-commit's fake builder reading the actor-record line from `rhiza guide builder`.

## statement

`rhiza guide ls` lists the builder, verifier and build-retry topics, guide install writes them as the rhiza-builder, rhiza-verifier and rhiza-build-retry skills beside the others, and the build topic (still installed as rhiza-build) is a short index that names the three role topics and says which one each role loads, so workflows that still say /rhiza-build or `rhiza guide build` keep working. cmd/rhiza/testdata/script/build-evidence-commit.txtar's fake builder, which greps the topic for the actor-record `git add` line, reads `rhiza guide builder` instead of `rhiza guide build`.

## Links

### Out

- depends-on [builder-topic](builder-topic.md)
- depends-on [verifier-topic](verifier-topic.md)
- depends-on [build-retry-topic](build-retry-topic.md)
- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)

### In

- depends-on [build-prompts-one-line](build-prompts-one-line.md)
- verifies [build-role-skills-installed-verified](../../../../../../evidence/build-role-skills-installed-verified.md)
- verifies [build-role-skills-installed-qa-1](../../../../../../evidence/build-role-skills-installed-qa-1.md)
- verifies [build-role-skills-installed-qa-2](../../../../../../evidence/build-role-skills-installed-qa-2.md)
- verifies [build-role-skills-installed-verified-5](../../../../../../evidence/build-role-skills-installed-verified-5.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
