# question-03f1e4cd-1

question · `01M79CTRMM1S9SKRRNTQTRA3VG`

## session

03f1e4cd-1fae-4504-9df6-01a1d2cbb3c5

## statement

Tiered verifier: what decides that a leaf is risky and gets the strongest verifier model? / Tiered verifier: which model verifies a low-risk leaf? / Tiered verifier: where are the risky paths and the two model ladders configured? / Tiered verifier: which verifier sessions does the tiering cover?

## Links

### In

- answers [tiered-verifier-rules](../decision/tiered-verifier-rules.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
