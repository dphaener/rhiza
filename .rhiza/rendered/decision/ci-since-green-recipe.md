# ci-since-green-recipe

decision · `01M7714KH6RX4DE5NJXB8HNN6A`

## quote

go with the recommendation on the placeholder finding, and a documented recipe

## quote_where

owner, 2026-11-12, to the orchestrator in chat, answering the plan gate questions of run 01M76STX

## statement

Other projects get the since-last-green CI rule as a documented recipe in the user docs, with an example changes job, in language-neutral words; rhiza ci install does not write it

## Links

### Out

- answers [ci-since-green-shipped-form](../question/ci-since-green-shipped-form.md)

### In

- decided-by [ci-since-green-for-projects](../capability/rhiza/integrity/ci-install/ci-since-green-for-projects.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
