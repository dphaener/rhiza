# land-conflict-script-named

decision · `01M7714JQCXVJAM178H1S333DB`

## quote

go with the recommendation on the placeholder finding, and a documented recipe

## quote_where

owner, 2026-11-12, to the orchestrator in chat, answering the plan gate questions of run 01M76STX

## statement

The acceptance of @land-conflict-to-fixer names its testscript run-ship-land-conflict (go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-land-conflict), and says the ship-change land step gains an on.fail that routes a land conflict (exit 3) to conflict-fix, not conflict-fix-landed, while any other land failure still fails land

## Links

### Out

- answers [review-01m76stxjb7cgga48pw43p487t-r3-f1](../question/review-01m76stxjb7cgga48pw43p487t-r3-f1.md)

### In

- decided-by [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
