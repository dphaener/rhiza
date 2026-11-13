# tiered-verifier-rules

decision · `01M79CX5B3P4T9BGG2V38PA0XS`

## reasoning

The change set's intent names the risky set (test infrastructure, hooks, CI, workflow files, migrations, retries) and the two outcomes (cheaper model for low-risk, strongest for risky). The owner accepted 'verifier: yes, tiered by risk' in the interview. In the code, every built-in verifier (internal/workflow/builtin/build.toml and review.toml) runs models = ["opus"], the builder ladder starts at sonnet, and models are picked per role by workflow.Role.Model in internal/run/leaves.go startLeaf. So the strongest model is opus and the cheaper one already trusted for this repo's builds is sonnet. Deciding risk from the commit's changed paths keeps the choice deterministic and testable at launch, with no extra agent call. A failed-then-retried leaf has already shown it is hard, so 'retries' is read as attempt, round or cycle past 1, plus a path glob for retry code. Failing safe to opus when the diff cannot be read keeps verification quality where it matters. Putting the new fields on the role, next to `models`, follows how ladders are configured today, keeps old workflows unchanged (not breaking), and avoids adding a rhiza.* setting. Covering review's fixed-leaf verifier too is consistent because it uses the same role and prompt, and that is where verifier wall time goes as well.

## statement

Tiered verifier rules. (1) Risk: when it launches a verifier, the conductor decides deterministically, with no agent judgement, that the leaf is risky if any file the commit under verification changes (the diff from the builder's or fixer's journalled start HEAD to $RHIZA_COMMIT) matches one of the role's risky path globs, or if the session is a retry: the leaf's attempt is past 1, it is in a fresh round, or (in review) it is past the review's first fix: fixes start in review cycle 2, so a cycle past 2 or a fix try past the first within its cycle. If the diff cannot be read, the leaf counts as risky. Everything else is low-risk. (2) A low-risk leaf is verified on sonnet; a risky one stays on opus, as every verifier is today. (3) Configuration lives on the verifier role in the workflow TOML, next to the ladder it already has: `models` stays the strong (risky) ladder, a new `low_risk_models` holds the cheap ladder, and a new `risky_paths` holds the globs. A role that sets neither new field behaves exactly as today. The built-in defaults are low_risk_models = ["sonnet"] and risky_paths covering test infrastructure (*_test.go, **/testdata/**, test helpers such as main_test.go), hooks (.githooks/**, **/hooks/**), CI (.github/**), workflow files (.rhiza/workflows/**, internal/workflow/builtin/**), migrations (**/migrations/**, **/migrate/**) and retry/backoff code (**/*retry*). A project changes them by overriding the workflow, as it already does for models; no new rhiza.* setting. (4) It covers every verifier session: build's per-leaf verifier and review's fixed-leaf verifier (step.review.verifier), since the rule is on the generic verifier role. Builders, fixers, checks and the critic keep their own ladders. Verifiers still write Evidence as themselves and build-check still requires a non-author verifier; the leaf-session entry journals the model chosen, so the tier is visible.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-03f1e4cd-1](../question/question-03f1e4cd-1.md)

### In

- decided-by [verifier-tier-fields](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-fields.md)
- decided-by [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)
- decided-by [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
