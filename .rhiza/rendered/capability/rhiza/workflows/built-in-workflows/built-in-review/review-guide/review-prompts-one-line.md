# review-prompts-one-line

capability · `01M79AH3SM4W5VEKB2SX2W6335` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/workflow/... passes with the review workflow test expecting the checks compliance, quality, qa and model, each check prompt naming its own topic and holding no session-done, test-scope or Checks: text, the critic and fixer prompts equal to the one-line forms, and a test reading .rhiza/workflows/docs-review.toml that asserts no scope check, its check and critic prompts equal the builtin ones, and its fixer prompt starts with the builtin fixer line and still names the humanize skill, doc-prose-humanized and doc-examples-plain-commands; go test -trimpath ./cmd/rhiza -run 'TestScripts/(review-workflow|review-checks-step|autonomous-review|review-summary)' passes with fakes that no longer expect a scope check; describe-golden is regenerated only if its output changed; make self-check passes.

## statement

The builtin review workflow has four checks (compliance, quality, qa, model; no scope), each check's prompt the one line `/rhiza-check-<check> Run the <check> check of change set {change} (cycle {cycle}): follow rhiza guide check-<check>.`, its critic prompt `/rhiza-critic Criticise the latest round of the {check} check of change set {change} (cycle {cycle}): follow rhiza guide critic.` and its fixer prompt `/rhiza-fixer Fix the findings below for change set {change} in review cycle {cycle}: follow rhiza guide fixer.`; Rhiza's .rhiza/workflows/docs-review.toml drops its scope check, uses the same check and critic lines byte for byte, and its fixer prompt is the builtin fixer line followed by its doc clause, kept as today: draft or revise every line of doc prose with the owner's humanize skill before committing, under constraints doc-prose-humanized and doc-examples-plain-commands (model writes shown as plain `rhiza add`, `link` and `set` commands; generated reference tables exempt).

## Links

### Out

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)

### In

- verifies [review-prompts-one-line-verified](../../../../../../evidence/review-prompts-one-line-verified.md)
- verifies [review-prompts-one-line-qa-1](../../../../../../evidence/review-prompts-one-line-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
