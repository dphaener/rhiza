## Severities

- `blocker`: behaviour broken, a build or test fails, or a leaf missing.
- `major`: an acceptance point unmet, or a real defect.
- `minor`: style, naming, docs or a non-blocking follow-up.
- `nit`: taste.

Per check: `compliance` leaf missing blocker, point unmet major; `quality` failing build or
test blocker, defect major, convention minor; `qa` broken behaviour blocker, rough edge
minor; `scope` outside the intent major; `model` code contradicted major, missing
constraint minor. Plan review: `missing-acceptance`, `contradiction`, `outside-intent`,
`uncovered-intent` and `oversized-leaf` major, the rest minor. `minor` and `nit` are never fixed.
