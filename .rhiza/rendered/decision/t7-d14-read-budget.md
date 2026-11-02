# t7-d14-read-budget

decision · `01M5SWSCG2YEVDXH7MTN70MBWJ`

## accepted_by

owner

## revised

revised by the orchestrator under the owner's authorization, when CI ran roadmap and ready at about 25.5 ms: budget = 50 ms (done point 4, as slice uses). Keep the 25 ms numbers as benchmarks; speeding up planGather is a follow-up.

## settled_by

Agent's option, accepted by the owner: Per-command budget test (Recommended): One test in make check: on the generated 50k model with a fresh index, run each read op through the api (show, ls first page, tree of a deep subtree, why, log, roadmap, ready) in as-is and to-be. Each median must be under 25 ms over 20 runs, with room for CI. Skipped under -short. what-if is exempt (D2). A whole-tree `tree` is benchmarked, not budgeted.

## statement

Budget test in `make check`: on the generated 50k model with a fresh index, each read op through the api (show, ls first page, tree of a deep subtree, why, log, roadmap, ready) in as-is and to-be; median < 25 ms over 20 runs; skipped under `-short`. what-if exempt; whole-tree `tree` benchmarked, not budgeted. **Revised in review by the owner's session:** the budget is 50 ms (milestone done point 4, as slice uses), since the CI runner needs ~25.5 ms for roadmap/ready; the 25 ms numbers stay as benchmarks; speeding up planGather is a follow-up.

## wording

agent of milestone-1 task 7 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-7](../source/plan-task-7.md)
- derived-from [task-interviews](../source/task-interviews.md)
- derived-from [founding-conversation](../source/founding-conversation.md)

### In

- relates-to [model-app-read-budget](model-app-read-budget.md) — applies t7-d14-read-budget's 50 ms tree-read budget to the model app's subtree read
- relates-to [whatif-overlay-reader](whatif-overlay-reader.md) — lifts t7-d14's what-if exemption: with the overlay reader, every budgeted read (show, ls first page, deep-subtree tree, why, log, roadmap, ready) is budgeted in what-if as well as as-is and to-be, under the same 50 ms median; the whole-tree tree stays benchmarked, not budgeted
- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS9YBW0KGX13FGMDHMKTG` Task 7's plan decisions (PR #10, merge 29e811d; PR #11, merge 3b5d1da), accepted by the owner
