# m29-test-root-per-script

decision · `01M6RYBXCNTMF66Z597WTFMN7R`

## statement

The testscript Setup makes each script's $WORK its own Rhiza root and writes $WORK/rhiza.toml registering a project named repo at $WORK/repo (directly, not through rhiza project add), so the existing scripts stay unedited; the acceptance checks that rhiza root show in $WORK/repo sees project repo, and $NOROOT stays for the outside-a-root cases.

## wording

Register work/repo: make each script's working folder its own root and register the repo folder inside it as a project named repo, keeping the existing scripts unedited; the agent's option, picked by the owner

## Links

### Out

- answers [review-01m6rmj6-r3-f1](../question/review-01m6rmj6-r3-f1.md)

### In

- decided-by [testscripts-run-under-root](../capability/rhiza/runtime/testscripts-run-under-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
