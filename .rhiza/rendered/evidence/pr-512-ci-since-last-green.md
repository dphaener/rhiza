# pr-512-ci-since-last-green

evidence · `01M791ZDNYTX41VSAW7Q61WP94` · actor [dphaener](../actor/dphaener.md)

## head

b03ae86fa26d3fea066f68b1f318d7d502f08066

## pr

https://github.com/dphaener/rhiza/pull/512

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #512 (https://github.com/dphaener/rhiza/pull/512) passed CI at b03ae86

## status

verified

## Links

### Out

- verifies [ci-model-only-since-green](../capability/rhiza/integrity/self-check/ci-model-only-since-green.md)
- verifies [ci-since-green-for-projects](../capability/rhiza/integrity/ci-install/ci-since-green-for-projects.md)
- verifies [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
