# umbrella-not-ready-verified

evidence · `01M6A20NMZ8RKV6DYMYVT5P6AQ` · actor [claude-code-b61e6f7b](../actor/claude-code-b61e6f7b.md)

## result

PASS

## statement

At 28dbfc9f: go test ./cmd/rhiza -run TestScripts/roadmap-umbrella passes; its new tail accepts umbrella U and milestone W, and ready --json shows the reason "umbrella: work is carved out of it" with next 01M5PWGY38TY4APFS7ZR0Q10WA (W, non-umbrella). go test ./internal/api -run 'TestUmbrellaNotReady|TestDaemon' passes; TestUmbrellaNotReady (in plan_test.go, package internal/api) serves an accepted change set to an umbrella intent, asserts planUmbrellaReason, and calls daemonCandidates with hasLeaves always true: no candidates, one SkipNotReady skip. The diff adds the reason in planCore for any served umbrella intent. make check passes with the RHIZA_ROOT/CHANGE/RUN/ACTOR env unset.

## Links

### Out

- verifies [umbrella-not-ready](../capability/rhiza/views/ready-next/umbrella-not-ready.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
