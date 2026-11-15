# fake-rex-no-dead-wait-verified

evidence · `01M7FY88REK69M0J1RWJW6B4KM` · actor [claude-code-ebd74164](../actor/claude-code-ebd74164.md)

## result

PASS

## statement

At leaf commit 87779b63e: go test -trimpath -p 2 -count=1 -v ./cmd/rhiza -run 'TestScripts/^(bugfix-workflow-autonomous|plan-workflow-autonomous-left|plan-workflow-autonomous|plan-review-workflow)$' exit 0, subtests 3.01, 5.04, 3.03, 2.51 s (all under 15 s; were 91, 64, 59, 58 s). The other 12 edited scripts (arbiter-agent-tabs, plan-milestone-backlog, plan-milestone-check, plan-review-visit-round, plan-workflow-ci-flake-red, plan-workflow-ci-flake, plan-workflow-ci-infra-wait, plan-workflow-merge-main-moved-often, plan-workflow-merge-main-moved, plan-workflow-open, plan-workflow-pr, plan-workflow) each pass run by name, 0.56-3.21 s. grep -l 'never got ready' lists 16 scripts, each with the started guard ([ -n "$started" ] || exit 1) set only past --, and the four slow ones plus arbiter-agent-tabs answer "block inspect" with a window id. The diff is additions only, inside the fake rex sections (no assertion lines changed), and touches nothing under internal/. rhiza ship generate and rhiza render exit 0.

## Links

### Out

- verifies [fake-rex-no-dead-wait](../capability/fake-rex-no-dead-wait.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
