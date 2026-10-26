# task-interviews

source · `01M5SWRGJXREMK1PV5Y23DKTF8`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Ship it, you merge it",
    "where": "Task 1, ship sign-off: typed answer"
  },
  {
    "by": "owner",
    "text": "yes, push it and open the PR, when CI is green let's merge it, remove the branches, and ship it",
    "where": "Task 2, after ship sign-off, 2026-10-24 13:55 UTC"
  },
  {
    "by": "owner",
    "text": "mark task 2 done and open a Rex tab for task 3",
    "where": "Task 2, 2026-10-24 13:58 UTC"
  },
  {
    "by": "owner",
    "text": "push it and open the PR, merge after CI is green, clean up branches, ship a binary.",
    "where": "Task 3, after ship sign-off, 2026-10-24 14:43 UTC"
  },
  {
    "by": "owner",
    "text": "yes mark it done, then open a Rex tab for task 4",
    "where": "Task 3, 2026-10-24 14:52 UTC"
  },
  {
    "by": "owner",
    "text": "yes, open the PR, wait for CI to be green, merge, then clean up the branches and ship a new binary.",
    "where": "Task 4, after ship sign-off, 2026-10-25 01:18 UTC"
  },
  {
    "by": "owner",
    "text": "yes, mark it done",
    "where": "Task 4, 2026-10-25 01:20 UTC"
  },
  {
    "by": "owner",
    "text": "rebase onto main and open the PR",
    "where": "Task 5, after ship sign-off, 2026-10-25 03:14 UTC"
  },
  {
    "by": "owner",
    "text": "push it and open the PR, wawit for CI to be green, merge it, clean up the branches, then ship a new binary.",
    "where": "Task 6, after ship sign-off, 2026-10-25 01:40 UTC"
  },
  {
    "by": "owner",
    "text": "remove the worktree",
    "where": "Task 6, 2026-10-25 01:43 UTC"
  },
  {
    "by": "owner",
    "text": "yes, ship it and open the PR, merge when CI is green, the clean up the branches, and ship a new binary.",
    "where": "Task 9, 2026-10-25 02:56 UTC"
  },
  {
    "by": "owner",
    "text": "mark task 9 done in milestone-1.md",
    "where": "Task 9, 2026-10-25 03:02 UTC"
  },
  {
    "by": "orchestrator (the owner's main session), not the owner",
    "text": "nto origin/main, run make check, push, open the PR, wait for CI to be green, and stop: the orchestrator merges and closes out. Do not merge.",
    "where": "Task 10, after ship sign-off, 2026-10-25 03:38 UTC; only the last 140 characters of a longer message arrived"
  },
  {
    "by": "orchestrator (the owner's main session), not the owner",
    "text": "Message from the orchestrator (the owner's main session; the owner authorized it to adjust plans): change D1's default view so reads and slice share one rule across tasks 7 and 8. Effective change set draft or proposed -> what-if of it; accepted -> to-be; none or landed -> as-is. Work happens on accepted change sets, so the main loop stays on the fast indexed to-be path; what-if stays out of the 50 ms budget. Update D1 in your plan, tell any running maker, and carry on. Task 7 is getting the same rule. No reply needed beyond noting it in the plan.",
    "where": "Task 8, D1 default-view amendment, 2026-10-25 03:32 UTC"
  },
  {
    "by": "orchestrator (the owner's main session), not the owner",
    "text": "Orchestrator: ship it. First rebase onto origin/main, regenerate describe-golden (it conflicts), make check, open the PR, stop at green CI.",
    "where": "Task 8, ship sign-off answer"
  },
  {
    "by": "orchestrator (the owner's main session), not the owner",
    "text": "Orchestrator: fixed on main in 3e17d7b, link/apply take change-set refs. Rebase onto origin/main, use genesis link, leave write_link.go.",
    "where": "Task 7, answer to the change-set link endpoint question"
  },
  {
    "by": "orchestrator (the owner's main session), not the owner",
    "text": "Orchestrator: budget = 50 ms (done point 4, as slice uses). Keep the 25 ms numbers as benchmarks; speeding up planGather is a follow-up.",
    "where": "Task 7, answer to the CI read budget question"
  }
]
```

## raw_sha256

4bcca6cf9b99fdff469c146181e5fd266c918b83fff797696bfe8f82b7bc1584

## summary

Every task-session interview for milestone 1 tasks 1-10, extracted by a script from the task sessions' Claude Code transcripts into an append-only record; the transcripts stay out of every repo. Each task's shipit agent asked questions with options it wrote, and the owner picked from them. Across the ten tasks the agents asked 197 questions. The owner accepted an agent option 193 times (at least 173 of them the option the agent marked recommended), typed one answer, and typed eleven more messages, all directions about shipping, merging, cleaning up branches and worktrees, and marking tasks done. Every design question went the agent's recommended way; the other picks were workflow choices (plan approval, fixing review findings before shipping, agent delegation, ship sign-off) and multi-select guardrail lists. Picked options are the agent's wording, accepted by the owner, so none is quoted here as the owner's words; the decisions they settled are the D-rows of each task's shipit plan. The orchestrator (the owner's main session) typed five answers into task sessions, recorded as not the owner's: a truncated close-out instruction to task 10, task 8's shared default-view amendment and its ship sign-off, and task 7's answers on change sets as link endpoints (already fixed on main) and on the CI read budget (50 ms, with speeding up roadmap and ready left as a follow-up). Topics by task: 1 code layout, lint, Go version, the testscript harness, CI, CLAUDE.md; 2 ULIDs, canonical JSON, strict reading, atomic writes, the log reader, the corpus; 3 operation revisions, lifecycle rules, strict bodies, node kinds, fractional keys; 4 the three views, causal order, field-level conflicts, history, conflict repair; 5 the pure-Go SQLite driver, index location, incremental catch-up, the generator, the budget test; 6 the api package, schemas, actors, references, apply, errors, lifecycle verbs; 7 the seven reads, default view, why, roadmap, ready, the read budget; 8 slice shape, the budget measure, constraint inheritance, cuts, the JSON shape; 9 render, check failures and warnings, the append-only base, ci install, the self-check; 10 status, source add, size caps, guide topics, install and the session-start hook. Nothing in the record is left open beyond the follow-ups it mentions.

## title

Task interviews, milestone 1 tasks 1-10

## Links

### In

- derived-from [t1-d1-version-usage](../decision/t1-d1-version-usage.md)
- derived-from [t1-d2-cmd-internal-layout](../decision/t1-d2-cmd-internal-layout.md)
- derived-from [t1-d3-gofmt-vet-lint](../decision/t1-d3-gofmt-vet-lint.md)
- derived-from [t1-d4-go-module-version](../decision/t1-d4-go-module-version.md)
- derived-from [t1-d5-testscript-in-process](../decision/t1-d5-testscript-in-process.md)
- derived-from [t1-d6-shipit-gitignored](../decision/t1-d6-shipit-gitignored.md)
- derived-from [t1-d7-branch-pr-ci](../decision/t1-d7-branch-pr-ci.md)
- derived-from [t1-d8-tdd-testscript](../decision/t1-d8-tdd-testscript.md)
- derived-from [t1-d9-version-stamping](../decision/t1-d9-version-stamping.md)
- derived-from [t1-d10-makefile-targets](../decision/t1-d10-makefile-targets.md)
- derived-from [t1-d11-short-claude-md](../decision/t1-d11-short-claude-md.md)
- derived-from [t1-d12-ci-one-job](../decision/t1-d12-ci-one-job.md)
- derived-from [t1-d13-task-1-out-list](../decision/t1-d13-task-1-out-list.md)
- derived-from [t1-d14-hermetic-git-scenario](../decision/t1-d14-hermetic-git-scenario.md)
- derived-from [t1-d15-task-1-acceptance](../decision/t1-d15-task-1-acceptance.md)
- derived-from [t2-d1-library-only](../decision/t2-d1-library-only.md)
- derived-from [t2-d2-ulid-ids](../decision/t2-d2-ulid-ids.md)
- derived-from [t2-d3-oklog-ulid](../decision/t2-d3-oklog-ulid.md)
- derived-from [t2-d4-tdd-rapid-fuzz](../decision/t2-d4-tdd-rapid-fuzz.md)
- derived-from [t2-d5-integer-numbers](../decision/t2-d5-integer-numbers.md)
- derived-from [t2-d6-minimal-escaping](../decision/t2-d6-minimal-escaping.md)
- derived-from [t2-d7-strict-reading](../decision/t2-d7-strict-reading.md)
- derived-from [t2-d8-opaque-body](../decision/t2-d8-opaque-body.md)
- derived-from [t2-d9-atomic-write](../decision/t2-d9-atomic-write.md)
- derived-from [t2-d10-collect-problems](../decision/t2-d10-collect-problems.md)
- derived-from [t2-d11-corpus-amendable](../decision/t2-d11-corpus-amendable.md)
- derived-from [t2-d12-task-2-out-list](../decision/t2-d12-task-2-out-list.md)
- derived-from [t2-d13-envelope-rules](../decision/t2-d13-envelope-rules.md)
- derived-from [t2-d14-four-packages](../decision/t2-d14-four-packages.md)
- derived-from [t2-d15-repo-json](../decision/t2-d15-repo-json.md)
- derived-from [t2-d16-no-time-check](../decision/t2-d16-no-time-check.md)
- derived-from [t2-d17-stray-files](../decision/t2-d17-stray-files.md)
- derived-from [t2-d18-corpus-location](../decision/t2-d18-corpus-location.md)
- derived-from [t2-d19-task-2-acceptance](../decision/t2-d19-task-2-acceptance.md)
- derived-from [t3-d1-pure-packages](../decision/t3-d1-pure-packages.md)
- derived-from [t3-d2-stable-op-ids](../decision/t3-d2-stable-op-ids.md)
- derived-from [t3-d3-revise-keeps-state](../decision/t3-d3-revise-keeps-state.md)
- derived-from [t3-d4-record-landing](../decision/t3-d4-record-landing.md)
- derived-from [t3-d5-deterministic-fold](../decision/t3-d5-deterministic-fold.md)
- derived-from [t3-d6-strict-bodies](../decision/t3-d6-strict-bodies.md)
- derived-from [t3-d7-free-fields-slug-handle](../decision/t3-d7-free-fields-slug-handle.md)
- derived-from [t3-d8-nine-node-kinds](../decision/t3-d8-nine-node-kinds.md)
- derived-from [t3-d9-fractional-keys](../decision/t3-d9-fractional-keys.md)
- derived-from [t3-d10-tdd-fuzz-properties](../decision/t3-d10-tdd-fuzz-properties.md)
- derived-from [t3-d11-task-3-guardrails](../decision/t3-d11-task-3-guardrails.md)
- derived-from [t3-d12-no-context-md](../decision/t3-d12-no-context-md.md)
- derived-from [t3-d13-three-packages](../decision/t3-d13-three-packages.md)
- derived-from [t3-d14-accept-reject-bodies](../decision/t3-d14-accept-reject-bodies.md)
- derived-from [t3-d15-replacement-rules](../decision/t3-d15-replacement-rules.md)
- derived-from [t3-d16-task-3-acceptance](../decision/t3-d16-task-3-acceptance.md)
- derived-from [t4-d1-pure-package](../decision/t4-d1-pure-package.md)
- derived-from [t4-d2-internal-model](../decision/t4-d2-internal-model.md)
- derived-from [t4-d3-causal-order](../decision/t4-d3-causal-order.md)
- derived-from [t4-d4-field-level-conflicts](../decision/t4-d4-field-level-conflicts.md)
- derived-from [t4-d5-conflicting-op-applies](../decision/t4-d5-conflicting-op-applies.md)
- derived-from [t4-d6-skip-broken-ops](../decision/t4-d6-skip-broken-ops.md)
- derived-from [t4-d7-one-engine-views](../decision/t4-d7-one-engine-views.md)
- derived-from [t4-d8-overlaps](../decision/t4-d8-overlaps.md)
- derived-from [t4-d9-duplicate-handles](../decision/t4-d9-duplicate-handles.md)
- derived-from [t4-d10-node-history](../decision/t4-d10-node-history.md)
- derived-from [t4-d11-tdd-properties-corpus](../decision/t4-d11-tdd-properties-corpus.md)
- derived-from [t4-d12-task-4-guardrails](../decision/t4-d12-task-4-guardrails.md)
- derived-from [t4-d13-task-4-acceptance](../decision/t4-d13-task-4-acceptance.md)
- derived-from [t4-d14-conflict-repair](../decision/t4-d14-conflict-repair.md)
- derived-from [t5-d1-index-scope](../decision/t5-d1-index-scope.md)
- derived-from [t5-d2-modernc-sqlite](../decision/t5-d2-modernc-sqlite.md)
- derived-from [t5-d3-index-location](../decision/t5-d3-index-location.md)
- derived-from [t5-d4-input-incremental](../decision/t5-d4-input-incremental.md)
- derived-from [t5-d5-stored-views](../decision/t5-d5-stored-views.md)
- derived-from [t5-d6-freshness-by-names](../decision/t5-d6-freshness-by-names.md)
- derived-from [t5-d7-generator-package](../decision/t5-d7-generator-package.md)
- derived-from [t5-d8-budget-test](../decision/t5-d8-budget-test.md)
- derived-from [t5-d9-index-failures](../decision/t5-d9-index-failures.md)
- derived-from [t5-d10-query-api](../decision/t5-d10-query-api.md)
- derived-from [t5-d11-generator-shape](../decision/t5-d11-generator-shape.md)
- derived-from [t5-d12-tests-after](../decision/t5-d12-tests-after.md)
- derived-from [t5-d13-budget-threshold](../decision/t5-d13-budget-threshold.md)
- derived-from [t5-d14-task-5-guardrails](../decision/t5-d14-task-5-guardrails.md)
- derived-from [t5-d15-no-adr](../decision/t5-d15-no-adr.md)
- derived-from [t6-d1-narrow-model-seam](../decision/t6-d1-narrow-model-seam.md)
- derived-from [t6-d2-frame-writes-init-scope](../decision/t6-d2-frame-writes-init-scope.md)
- derived-from [t6-d3-api-registry](../decision/t6-d3-api-registry.md)
- derived-from [t6-d4-reflected-schemas](../decision/t6-d4-reflected-schemas.md)
- derived-from [t6-d5-actor-lookup](../decision/t6-d5-actor-lookup.md)
- derived-from [t6-d6-current-change-file](../decision/t6-d6-current-change-file.md)
- derived-from [t6-d7-tdd-testscript](../decision/t6-d7-tdd-testscript.md)
- derived-from [t6-d8-ref-resolution](../decision/t6-d8-ref-resolution.md)
- derived-from [t6-d9-revise-existing-op](../decision/t6-d9-revise-existing-op.md)
- derived-from [t6-d10-apply-file-format](../decision/t6-d10-apply-file-format.md)
- derived-from [t6-d11-interim-model-unexported](../decision/t6-d11-interim-model-unexported.md)
- derived-from [t6-d12-error-shape](../decision/t6-d12-error-shape.md)
- derived-from [t6-d13-lifecycle-verbs](../decision/t6-d13-lifecycle-verbs.md)
- derived-from [t6-d14-guardrails](../decision/t6-d14-guardrails.md)
- derived-from [t6-d15-done](../decision/t6-d15-done.md)
- derived-from [t7-d1-reads-scope](../decision/t7-d1-reads-scope.md)
- derived-from [t7-d2-default-view](../decision/t7-d2-default-view.md)
- derived-from [t7-d3-reads-use-index](../decision/t7-d3-reads-use-index.md)
- derived-from [t7-d4-shared-read-plumbing](../decision/t7-d4-shared-read-plumbing.md)
- derived-from [t7-d5-why-walk](../decision/t7-d5-why-walk.md)
- derived-from [t7-d6-log-history](../decision/t7-d6-log-history.md)
- derived-from [t7-d7-change-set-links](../decision/t7-d7-change-set-links.md)
- derived-from [t7-d8-ready-rule](../decision/t7-d8-ready-rule.md)
- derived-from [t7-d9-roadmap-grouping](../decision/t7-d9-roadmap-grouping.md)
- derived-from [t7-d10-tree-shape](../decision/t7-d10-tree-shape.md)
- derived-from [t7-d11-ls-filters-paging](../decision/t7-d11-ls-filters-paging.md)
- derived-from [t7-d12-show-shape](../decision/t7-d12-show-shape.md)
- derived-from [t7-d13-tests](../decision/t7-d13-tests.md)
- derived-from [t7-d14-read-budget](../decision/t7-d14-read-budget.md)
- derived-from [t7-d15-guardrails](../decision/t7-d15-guardrails.md)
- derived-from [t7-d16-finish](../decision/t7-d16-finish.md)
- derived-from [t7-d17-plan-reads-to-be](../decision/t7-d17-plan-reads-to-be.md)
- derived-from [t8-d1-slice-shape](../decision/t8-d1-slice-shape.md)
- derived-from [t8-d2-budget-stored-views](../decision/t8-d2-budget-stored-views.md)
- derived-from [t8-d3-budget-printed-bytes](../decision/t8-d3-budget-printed-bytes.md)
- derived-from [t8-d4-token-estimate](../decision/t8-d4-token-estimate.md)
- derived-from [t8-d5-constraints-via-links](../decision/t8-d5-constraints-via-links.md)
- derived-from [t8-d6-short-ids](../decision/t8-d6-short-ids.md)
- derived-from [t8-d7-slice-package](../decision/t8-d7-slice-package.md)
- derived-from [t8-d8-change-slice-focus](../decision/t8-d8-change-slice-focus.md)
- derived-from [t8-d9-item-fields](../decision/t8-d9-item-fields.md)
- derived-from [t8-d10-cut-order](../decision/t8-d10-cut-order.md)
- derived-from [t8-d11-tests](../decision/t8-d11-tests.md)
- derived-from [t8-d12-guardrails](../decision/t8-d12-guardrails.md)
- derived-from [t8-d13-current-decisions](../decision/t8-d13-current-decisions.md)
- derived-from [t8-d14-neighbours-in-flight](../decision/t8-d14-neighbours-in-flight.md)
- derived-from [t8-d15-over-budget](../decision/t8-d15-over-budget.md)
- derived-from [t8-d16-json-shape](../decision/t8-d16-json-shape.md)
- derived-from [t8-d17-proposed-change](../decision/t8-d17-proposed-change.md)
- derived-from [t9-d1-render-as-is](../decision/t9-d1-render-as-is.md)
- derived-from [t9-d2-render-explicit](../decision/t9-d2-render-explicit.md)
- derived-from [t9-d3-check-fail-warn](../decision/t9-d3-check-fail-warn.md)
- derived-from [t9-d4-check-base](../decision/t9-d4-check-base.md)
- derived-from [t9-d5-ci-install-files](../decision/t9-d5-ci-install-files.md)
- derived-from [t9-d6-self-check](../decision/t9-d6-self-check.md)
- derived-from [t9-d7-check-report](../decision/t9-d7-check-report.md)
- derived-from [t9-d8-tests](../decision/t9-d8-tests.md)
- derived-from [t9-d9-node-page](../decision/t9-d9-node-page.md)
- derived-from [t9-d10-endpoint-table](../decision/t9-d10-endpoint-table.md)
- derived-from [t9-d11-stale-render](../decision/t9-d11-stale-render.md)
- derived-from [t9-d12-guardrails](../decision/t9-d12-guardrails.md)
- derived-from [t9-d13-defaults](../decision/t9-d13-defaults.md)
- derived-from [t9-d14-acceptance](../decision/t9-d14-acceptance.md)
- derived-from [t10-d1-minimal-status](../decision/t10-d1-minimal-status.md)
- derived-from [t10-d2-source-add-shape](../decision/t10-d2-source-add-shape.md)
- derived-from [t10-d3-size-cap](../decision/t10-d3-size-cap.md)
- derived-from [t10-d4-tests](../decision/t10-d4-tests.md)
- derived-from [t10-d5-topics](../decision/t10-d5-topics.md)
- derived-from [t10-d6-guide-command](../decision/t10-d6-guide-command.md)
- derived-from [t10-d7-full-loop-allowlist](../decision/t10-d7-full-loop-allowlist.md)
- derived-from [t10-d8-status-hook](../decision/t10-d8-status-hook.md)
- derived-from [t10-d9-guardrails](../decision/t10-d9-guardrails.md)
- derived-from [t10-d10-done](../decision/t10-d10-done.md)
- derived-from [t10-d11-skills-dir-plugin](../decision/t10-d11-skills-dir-plugin.md)
- derived-from [t10-d12-reinstall-rule](../decision/t10-d12-reinstall-rule.md)
- derived-from [t10-d13-hook-exec-form](../decision/t10-d13-hook-exec-form.md)
- derived-from [area-gap](../question/area-gap.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
