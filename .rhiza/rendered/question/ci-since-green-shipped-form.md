# ci-since-green-shipped-form

question · `01M76XZFWK2JJH0EV5MZQ9FAS9`

## statement

rhiza ci install writes only a rhiza check workflow (.github/workflows/rhiza.yml), never a project's test jobs, so the since-last-green rule has no shipped CI template to go in. How should other projects get it: a documented recipe in the user docs with an example changes job, an optional changes job that rhiza ci install writes for the project's own jobs to gate on, or Rhiza's repo only?

## Links

### Out

- blocks [ci-since-green-for-projects](../capability/rhiza/integrity/ci-install/ci-since-green-for-projects.md)

### In

- answers [ci-since-green-recipe](../decision/ci-since-green-recipe.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
