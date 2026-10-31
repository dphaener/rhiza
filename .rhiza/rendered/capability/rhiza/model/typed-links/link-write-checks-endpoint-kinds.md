# link-write-checks-endpoint-kinds

capability · `01M6857PTDC3M6J8MAWWXFF0ZS` · parent [typed-links](../typed-links.md)

## acceptance

link.txtar checks rhiza link store serves log --json, with --record and with --dry-run, each exits 2 with code usage and the message link type "serves" does not go from capability to capability; a change set serving a capability is refused and a change set serving an intent succeeds; an apply.txtar case with one bad link item refuses the whole apply and adds no entry; a Go table test shows buildLink agrees with check.Allowed for every type and end kinds; existing scripts that wrote bad links are changed to valid ones

## statement

rhiza link, apply link items and --record links refuse a link whose type does not go from its from-end's kind to its to-end's kind, by the endpoint table rhiza check uses, as a usage error naming both kinds and writing nothing

## Links

### Out

- decided-by [link-kinds-refused-at-write](../../../../decision/link-kinds-refused-at-write.md)

### In

- verifies [link-write-checks-endpoint-kinds-verified](../../../../evidence/link-write-checks-endpoint-kinds-verified.md)
- verifies [link-write-checks-endpoint-kinds-qa-1](../../../../evidence/link-write-checks-endpoint-kinds-qa-1.md)
- verifies [pr-124-genesis-link-checks-endpoint](../../../../evidence/pr-124-genesis-link-checks-endpoint.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
