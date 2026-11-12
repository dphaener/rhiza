# note-01m76btstcnfv7

note · `01M76BTSTCNFV75E7CAME7TF1S`

## statement

Intent narrowed (plan review f3): release tooling is for this repo only (@answer-41425f42-20), so the marks are notes on change sets read by tools/release (@release-marks), not fields of rhiza requests. A severe fix is marked by a severe: note on its change set, not a severity on the bug request, and the rhiza build a bug was filed against is not recorded on the request. The soak gate (@@release-gates) instead counts any bug request filed since the root daemon first ran the candidate commit, which is conservative: every such bug was filed against that commit or a later one.

## Links

### Out

- discusses change set `01M76B077D23CJ28C690Y166VY` (Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
