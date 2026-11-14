# review-01m7ah6rayrzax6z6r63q0m2qk-r6-f5

question · `01M7BZ7ENTZ4H1WGF0YX6DE5P1`

## statement

Plan review round 6 left finding f5 (contradiction on builtin-conductor-default): Main's CI is one step running rhiza ship check --main. If the merge commit's tree differs from the tested head (main moved between ship merge's check and GitHub's merge), ship check --main dispatches CI on the merge commit and has to wait for it to settle (internal/api/ship_main.go shipCheckMainDispatch). Classic does this with a wait step, wait-main (until = rhiza ship check --main --settled), before main-green. A single command step checks once and goes red while the dispatched CI is still pending, which the plan treats as main being red. The critic confirmed it: Verified: origin/main shipCheckMainDispatch (internal/api/ship_main.go lines 101-106) returns Passed=false while CI is not settled, and ShipCheckResponse.ExitCode (ship.go:320-328) maps that to exit 1. Classic waits with step wait-main (kind wait, until = rhiza ship check --main --settled, builtin/ship-change.toml:468-471) before main-green. The leaf's single main-CI step with no on.fail route would stop the run as 'main is red' while dispatched CI is still pending. Major stands. How should it be settled?

## Links

### Out

- blocks [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)

### In

- answers [main-ci-one-wait-step](../decision/main-ci-one-wait-step.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
