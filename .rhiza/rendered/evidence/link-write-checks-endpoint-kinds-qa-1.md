# link-write-checks-endpoint-kinds-qa-1

evidence · `01M687R4G1J8WAM4Q072AQTMNY` · actor [claude-code-d48993a5](../actor/claude-code-d48993a5.md)

## result

PASS

## statement

QA at 9c8d4bf with a scratch build in a fresh repo: in an open change set rhiza link store serves log (capability to capability) exits 2 with code usage and message link type "serves" does not go from capability to capability under --json, plain, --dry-run and --record, writing nothing; link fast verifies store gives intent to capability, link c1 exempts log gives constraint to capability, both exit 2; link store serves fast and link store exempts c1 --note x --decision d1 succeed; linking the change set serves store is refused (change set to capability) and serves fast succeeds; rhiza apply with an add decision plus one bad link exits 2 naming writes[1], the log entry count stays 9 and later-link is not added.

## Links

### Out

- verifies [link-write-checks-endpoint-kinds](../capability/rhiza/model/typed-links/link-write-checks-endpoint-kinds.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
