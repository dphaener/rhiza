# t9-d5-ci-install-files

decision · `01M5SWSM1DZ14RX3E1NV2EAND4`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Workflow + gitattributes (Recommended): Writes .github/workflows/genesis.yml (checkout with fetch-depth 0, `go install github.com/dphaener/genesis/cmd/genesis@<this binary's version>`, `genesis check --base origin/$GITHUB_BASE_REF` on PRs) and adds `.genesis/rendered/** linguist-generated=true` to .gitattributes. Idempotent; an existing different workflow is a conflict unless --force.

## statement

`ci install` writes `.github/workflows/genesis.yml` and adds `.genesis/rendered/** linguist-generated=true` to `.gitattributes`; idempotent; a differing existing workflow is a conflict unless `--force`

## wording

agent of milestone-1 task 9 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-9](../source/plan-task-9.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [ci-install](../capability/genesis/integrity/ci-install.md)

## History

- `01M5SWSK1HH2VMB5RJ89DJ6YBX` Task 9's plan decisions (PR #6, merge 8355823), accepted by the owner
