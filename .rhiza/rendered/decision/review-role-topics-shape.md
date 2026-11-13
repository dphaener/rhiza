# review-role-topics-shape

decision · `01M79ADQKB796M6D1PD488C1M9`

## reasoning

Read the change set intent (nine named roles, six-section skeleton, isolation, scope leaves review.toml), prompt-gap.md sections 2.5-2.8 and 5 (CS1 names the topics check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix, conflict-fix; topics are skills installed by guide_install.go; builtin prompts become `/rhiza-<role> <vars>`; shared blocks include severity defs), the owner's interview picks (role topics in three change sets; fold scope into compliance, decision 01M76Y5M), and the sibling Build role topics arbiter decision (01M7971N), whose shared-block markers, rhiza-<topic> skills, one-line prompts and keep-the-old-topic-as-index pattern this follows so the three role change sets read alike. review.toml already gives checks a throwaway checkout told only change set, check and cycle; prompt-gap 2.5.1 says the gap is that nothing stops a check reading `rhiza run status` and anchoring on the build's findings, so a MUST NOT in the topics closes it without Go work, and tool guards are CS3's. Scope criteria follow afk's compliance check (under-built, scope creep, unaccounted files) and keep today's grading (outside intent major, leaf missing blocker). grep shows no Go code names a `scope` check, so check names are free-form; AGENTS.md calls breaking only removed or changed commands, flags, --json shapes, settings or .rhiza/ formats, and a builtin workflow's default check list is none of those, especially with the review index routing `scope`. The fixer steps were mapped from internal/workflow/builtin/ship-change.toml: gate-fix/gate-fix-merge share one prompt, conflict-fix/conflict-fix-landed share one, and model-conflict-fix is the other conflict kind, so one conflict-fix topic with two sections avoids a tenth topic the intent does not list.

## statement

Names and delivery: nine topics, `check-compliance`, `check-quality`, `check-qa`, `check-model`, `critic`, `fixer`, `ci-fix`, `gate-fix` and `conflict-fix`, each on the six-section skeleton (TASK, EXPECTED OUTCOME, READ FIRST, MUST DO, MUST NOT DO, CONTEXT), each installed by guide.install as a skill named rhiza-<topic> (rhiza-check-compliance ... rhiza-conflict-fix), and each reusing the build role topics' shared-block files under internal/guide/text/shared/ through `<!-- shared: <block> -->` lines, plus a new shared `severity` block (the definitions and the per-check mapping now in review.md) included by every check, the critic and the fixers. The builtin review.toml prompts shrink to one line carrying the TASK variables and both routes, e.g. `/rhiza-check-compliance Run the compliance check of change set {change} (cycle {cycle}): follow rhiza guide check-compliance.`; since one [step.review.check] prompt serves one check, each check's prompt names its own topic. The `review` topic and the rhiza-review skill stay, cut to a short index naming the nine topics and saying which to load (a check named `scope` loads check-compliance), so custom workflows still saying /rhiza-review or `rhiza guide review` keep working. Isolation: the builtin check and critic prompts keep carrying only change set, check and cycle, and checks and critics keep their throwaway checkout; the check-* and critic topics' MUST NOT DO says they know nothing of how the work was built and must not seek it: never read the build run's journal, findings, run notes or session transcripts (no `rhiza run status`/`run show` of the build run, no builder retries); judge only the contract (intent, leaves, acceptance, constraints from the slice), the diff from the merge-base and the model's Evidence. A critic reads only its own check's last round of this review run. No Go guard or prompt-assembly change; CS3 owns tool guards. Scope: the builtin review.toml loses its `scope` check (four checks per cycle), and Rhiza's own .rhiza/workflows/docs-review.toml drops it too. Its criteria move into check-compliance, graded beside each leaf's acceptance: under-built (a leaf missing: blocker; an acceptance point unmet: major), scope creep (behaviour outside the intent and the leaves: major) and unaccounted files (a changed file no leaf or the intent accounts for: major if it changes behaviour, minor otherwise). Not breaking, no `breaking:` note: no CLI command, flag, --json shape, rhiza.* setting or .rhiza/ file format changes; check names are free-form, so a repo workflow still listing `scope` validates and runs through the review index. Fixer steps: the review step's fixer uses `fixer`; ship-change's gate-fix and gate-fix-merge use `gate-fix`; ci-fix uses `ci-fix`; conflict-fix, conflict-fix-landed and model-conflict-fix use `conflict-fix` (one topic with a git-merge section and a model-conflict section, `rhiza set --rebase`). Each of those prompts shrinks to one line naming its topic and its TASK variables ({change}, {run}, and for model-conflict-fix that it is the model conflict), in both internal/workflow/builtin/ship-change.toml and Rhiza's .rhiza/workflows/ship-change.toml and docs.toml; everything else in today's prompts (reading the gate log, reproduce-before-fix, flake quarantine with a bug, trailers, staging .rhiza, the clash question, never push/rebase/force, session done as its own command) moves into the topics. plan.toml's prepush-fix is left to Plan role topics (01M76XRP).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-acb282c4-1](../question/question-acb282c4-1.md)
- derived-from [review-role-topics-interview](../source/review-role-topics-interview.md)

### In

- decided-by [severity-shared-block](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-shared-block.md)
- decided-by [check-compliance-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-compliance-topic.md)
- decided-by [check-quality-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-quality-topic.md)
- decided-by [check-qa-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-qa-topic.md)
- decided-by [check-model-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-model-topic.md)
- decided-by [critic-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/critic-topic.md)
- decided-by [fixer-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-topic.md)
- decided-by [ci-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/ci-fix-topic.md)
- decided-by [gate-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/gate-fix-topic.md)
- decided-by [conflict-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/conflict-fix-topic.md)
- decided-by [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)
- decided-by [review-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-prompts-one-line.md)
- decided-by [ship-fix-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/ship-fix-prompts-one-line.md)
- decided-by [review-guide](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide.md)
- decided-by [review-workflow-file](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-workflow-file.md)
- decided-by [severity-contract](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/severity-contract.md)
- decided-by [built-in-review](../capability/rhiza/workflows/built-in-workflows/built-in-review.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
