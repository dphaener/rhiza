## Severities

- `blocker`: behaviour broken, build or test failing, leaf missing.
- `major`: an acceptance point unmet, real defect.
- `minor`: style, naming, docs, follow-up.
- `nit`: taste.

Per check: `compliance` leaf missing blocker, point unmet or out of scope major;
`quality` failing build or test blocker, defect major, convention minor; `qa` breakage
blocker, rough edge minor; `model` code contradicted major, missing constraint minor.
Plan review: `missing-acceptance`, `contradiction`, `outside-intent`,
`uncovered-intent`, `oversized-leaf` and `no-failure-case` major, the rest minor.
`minor` and `nit` are never fixed.
