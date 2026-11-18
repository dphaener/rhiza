# api-tests-half-cpu-accept-at-60

decision · `01M7P7A06RNHA3KFE8P1CZCJ77`

## reasoning

The half-CPU line is a per-leaf proxy, not a milestone target; the intent's 'never lower a target' and m39-time-profile-first bind milestone targets, budget thresholds, rapid counts and tests, and this ruling touches none of them. The builder's two attempts (1a23840c3, eee179680), confirmed by the verifier at 198.84 s, cut the index setup the profile names (index.Open 151.85 s to 76.33 s of samples, catch-up 131.70 s to 65.46 s) through per-test copies of a prebuilt index, which is what m39-api-index-per-test-copy decided. Each trade-off fails on its own terms: sharing one index across runs of the binary is a cross-run shared index file and a non-hermetic cache, against the leaf's must_not and api-tests-on-index; dropping the fold side of TestWritesThroughIndexMatchFold removes the comparison the test exists for, which must_not forbids as shortening a test; baking the budgets' setup writes into the generated model changes the state the budgets measure, sits in budget tests outside the leaf's files (which name snapshot-template and snapshot-copy tests), and by the question's own numbers saves at most about 14 s (4.8 + 2.7 + 2.5 + 4.1), leaving about 180 s, still above 160.6 s. So another attempt could not pass the line within must_not, and failing the leaf would block cold-run-ship-scripts-cut, which carries the wall-time check that actually measures the milestone. Accepting the real 40% cut and leaving the milestone's 240 s judged downstream keeps every target and the tests' meaning intact

## statement

Accept leaf api-tests-copy-prebuilt-index at the measured cut (after user+sys 193.69 s, about 60% of the before 321.25 s; the verifier's 198.84 s) in place of its half-CPU line, provided the rest of its acceptance holds as written: the before/after Checks: lines and pprof index-setup lines, two -shuffle=on passes with different seeds, a test count no lower than before, and the fixture_test.go failure case (an apply on a snapshot copy read back through its copied index, the template's index file unchanged in size and mtime). None of the three named trade-offs is allowed: no setup writes baked into the budgets' generated model, no index shared across runs of the test binary, no dropping or narrowing the fold side of TestWritesThroughIndexMatchFold. The milestone targets are not lowered: the wall-time acceptance of cold-run-ship-scripts-cut (m39-cold-leaf-wall-proportion) and the cold make test of m39-cold still judge the milestone, and whatever internal/api still costs after that is planned as its own change set

## wording

(assumed, would have escalated)

## Links

### Out

- answers [api-tests-half-cpu-unreachable](../question/api-tests-half-cpu-unreachable.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
