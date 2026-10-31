# roadmap-umbrella-group-qa-1

evidence · `01M6A3103NECWM70QK0X3D4EDV` · actor [claude-code-36a5beca](../actor/claude-code-36a5beca.md)

## result

PASS

## statement

Built 3c614aeb to a scratch binary and, in a scratch git repo with RHIZA_* unset, added intents m1 (milestone=1), umb and umb2 (umbrella=yes), and proposed change sets U (serves umb), V (serves m1 and umb), U2 (serves umb2), W (serves m1) and X (serves nothing). rhiza roadmap --json gave groups milestone m1 [W], backlog [X], umbrella umb [U, V], umbrella umb2 [U2] in that order. V appears only under umb, not under m1. The human rhiza roadmap printed '## m1', '## Backlog', '## umb', '## umb2', and rhiza api call roadmap '{}' gave the same bytes as roadmap --json. go test -run TestScripts/(roadmap-umbrella|describe-golden) passes.

## Links

### Out

- verifies [roadmap-umbrella-group](../capability/rhiza/views/roadmap/roadmap-umbrella-group.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
