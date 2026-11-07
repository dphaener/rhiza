# decision-ci-triage-runners

decision · `01M6TBCFT3VTPXC20R9C8F9TQT`

## reasoning

The change set's slice and why list no intent, constraint or decision, so I read the model's intents: m27-nebari-on-rhiza ('Nebari's daemon builds and ships a Nebari change set on its own with Rhiza's workflows') is exactly what nebari#72's failed run blocks, since a Rails project's red CI never reaches ci-fix. In internal/api/ship.go, parseFailedTests matches only '--- FAIL: (\S+)' and shipReadFailures turns each name into a ShipFinding whose reason says 'rerun it alone and fix it'; ship-change.toml's ci-fix prompt tells the fixer to re-run only those tests. So the name must be something a fixer can rerun without guessing the runner: the runner's own printed rerun line is that, and it cannot be mistaken for a Go name. The change set's intent asks for RSpec plus 'other common runners'; I chose those that print one unambiguous per-test rerun line (RSpec, Minitest as Rails uses it, pytest) and left out count lines and file-level markers, which would turn infra failures into fake test findings and send them to ci-fix instead of ending the run. Go behaviour, including keeping only subtest leaves, stays as it is so existing ship tests hold.

## statement

Change set 01M6TB6N serves intent m27-nebari-on-rhiza (01M6J303FW4PYT5E5CE39QNHJ8). ship check's parseFailedTests keeps Go's '--- FAIL: <name>' lines (leaves only) and also reads, from gh run view --log-failed, each runner's own rerun line for a failed test, matched at the start of the log text after gh's job/step/timestamp prefix: RSpec's 'rspec ./<path>:<line>' lines (the 'Failed examples:' block), Rails Minitest's 'rails test <path>:<line>' / 'bin/rails test <path>:<line>' lines, and pytest's 'FAILED <path>::<test>' summary lines; nothing broader (no 'N failures' counts, no Jest file-level 'FAIL <file>'), so a log without one of these lines still ends as an infra or unparsable failure. A non-Go failing test is named by the rerun command its runner printed, trimmed of any trailing comment: 'rspec ./spec/views/tasks/list_spec.rb:290', 'bin/rails test test/models/x_test.rb:12', 'pytest tests/test_x.py::test_y'; Go tests keep their bare names (TestScripts/x). That name goes unchanged into ship check's tests and each finding's Test, deduplicated, and the tests field's description says it names failing tests of Go, RSpec, Minitest and pytest.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-bfd0cb38-1](../question/question-bfd0cb38-1.md)
- derived-from [ci-triage-runners-plan](../source/ci-triage-runners-plan.md)

### In

- decided-by [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests
