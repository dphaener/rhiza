# lint-hits-feed-model-check

decision · `01M7AZNSY3NEDH37CH6NB8XDH5`

## reasoning

The intent asks for deterministic lints that save the model check's cycle-1 failures (37-42%, mostly mechanical drift). The builtin review workflow's code step has no on.fail, so making lint hits fail code-check would end the run failed. A failed check is fixed and re-run in the next cycle (internal/run/checks.go checkVerdict; fixes start in cycle 2), so a hit the model check confirms as a major in cycle 1 still forces cycle 2: the hits must be fixed before the checks step runs. A command step that fails to a fixer agent step, which always passes on, does that with the existing step kinds and the existing fixer role, and false positives cost only an unfixed line. Scoping the stale-reference lint to the slice matches how review works today and keeps an unrelated change set from failing over others' drift; the branch's own deleted paths are the one deterministic, branch-scoped signal for a capability outside the slice whose code this branch removed. For Decisions, the model has decided-by links (node -> Decision) and builders' Leaf: trailers, both readable without judgement, so they give a deterministic test; arbiter Decisions like this one link only answers, hence the exemption. Revised in plan review round 1 (findings f1, f2).

## statement

(1) A lint hit never fails `rhiza run code-check`: the code step still fails only on a rhiza check failure or a leaf lacking verifier Evidence, as today. The lints run inside code-check and their hits are journalled on the code-check entry (shown by rhiza run status). The hits are fixed before the checks step runs, so the cycle-1 checks already see the fixed branch: the builtin review routes the code step's pass to a `lints` command step (`rhiza run lint-check`), which passes straight to the checks step when the latest code-check entry has no hits and otherwise fails to a `lint-fix` agent step; that step runs one fixer told the hits, who fixes each hit it agrees with (stale statement or acceptance text, an unretired capability, an unapplied Decision) in one commit trailed Review:, leaves the rest listed as unfixed, and always passes on to the checks step. A confirmed model-check major forces a second cycle, so handing hits to the model check as cycle-1 leads would not save one; the model check gets no leads. No on.fail is added to the code step. (2) The stale-reference lint checks the run's change set's slice only: its leaves and the capabilities and other nodes it adds or changes, plus their acceptance. It does not lint the rest of the model; drift elsewhere belongs to other change sets and would only add noise. A live capability outside the slice is flagged only when this branch removed its code: every path it references is among the paths deleted in merge-base..HEAD or otherwise absent at head, and at least one was deleted on the branch. (3) A Decision is "never applied" when the change set adds or changes it and either no live node links to it with decided-by, or every node that does is a leaf of this change set with no commit carrying its Leaf: trailer on the branch. A Decision whose only link is answers on a question is exempt from the first clause (an answered question binds no node by itself). The lint is mechanical and makes no judgement of code meaning; that stays with the fixer and the model check.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c0130112-1](../question/question-c0130112-1.md)
- derived-from [model-lints-plan-interview](../source/model-lints-plan-interview.md)

### In

- decided-by [modellint-stale-refs](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-stale-refs.md)
- decided-by [modellint-unapplied-decisions](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-unapplied-decisions.md)
- decided-by [code-check-runs-lints](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)
- decided-by [review-lint-fix-route](../capability/review-lint-fix-route.md)
- decided-by [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)
- decided-by [modellint-branch-removed-code](../capability/modellint-branch-removed-code.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
