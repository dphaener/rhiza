# umbrella-not-ready-qa-1

evidence · `01M6A311HERB0VNRKAPVBV0WAG` · actor [claude-code-36a5beca](../actor/claude-code-36a5beca.md)

## result

PASS

## statement

In the same scratch repo I accepted U, W and U2. rhiza ready --json marked U and U2 not ready with reason 'umbrella: work is carved out of it', marked V with that reason plus 'not accepted: proposed', and W was ready with next = W. The human rhiza ready printed 'next: ... W milestone' and listed the umbrella reason under each umbrella change set. go test ./internal/api passes, including TestUmbrellaNotReady.

## Links

### Out

- verifies [umbrella-not-ready](../capability/rhiza/views/ready-next/umbrella-not-ready.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
