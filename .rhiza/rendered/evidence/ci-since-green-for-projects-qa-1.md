# ci-since-green-for-projects-qa-1

evidence · `01M790GBD8J7T5QSJR1QRE6GZ4` · actor [claude-code-c8e6c442](../actor/claude-code-c8e6c442.md)

## result

PASS

## statement

QA at 4ed003163: read docs/ci-since-last-green.md and its README link; the recipe's rule matches the script's behaviour exercised by hand (no green, non-ancestor green, .rhiza/-only diff vs code diff), its aggregate job ci-passed does not clash with the check job rhiza ci install writes, and its green test requires every check run of that name to succeed. go test -trimpath ./internal/api -run 'Docs|CI' passes.

## Links

### Out

- verifies [ci-since-green-for-projects](../capability/rhiza/integrity/ci-install/ci-since-green-for-projects.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
