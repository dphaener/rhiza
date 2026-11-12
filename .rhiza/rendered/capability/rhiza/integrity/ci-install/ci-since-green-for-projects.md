# ci-since-green-for-projects

capability · `01M76XZFVQM9EKM6YN432SYK2V` · parent [ci-install](../ci-install.md)

## acceptance

Set by the answer to question ci-since-green-shipped-form, as tests of the chosen form (a docs page check in go test -trimpath ./internal/api -run Docs, or a ci.install test in go test -trimpath ./internal/api -run CIInstall).

## statement

Projects other than Rhiza get the since-last-green rule in the form the owner chooses, in language-neutral words.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../decision/efficient-cycle-milestone-shape.md)
- decided-by [ci-since-green-recipe](../../../../decision/ci-since-green-recipe.md)

### In

- blocks [ci-since-green-shipped-form](../../../../question/ci-since-green-shipped-form.md)
- verifies [ci-since-green-for-projects-verified](../../../../evidence/ci-since-green-for-projects-verified.md)
- verifies [ci-since-green-for-projects-qa-1](../../../../evidence/ci-since-green-for-projects-qa-1.md)
- verifies [pr-512-ci-since-last-green](../../../../evidence/pr-512-ci-since-last-green.md)

## History

- `01M76XSFD0JZ9424EWBSQA3K8X` CI on a pull request runs only the model checks when nothing outside .rhiza/ changed since the pull request's last commit whose CI passed, so ship land's push and merges of main that bring only .rhiza/ changes cost a model check instead of a full CI round. Evidence: wait-landed is a full CI round in every run, 5.9m median (perf report rec 1).
