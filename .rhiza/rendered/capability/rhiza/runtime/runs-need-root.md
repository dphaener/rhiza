# runs-need-root

capability · `01M6RN852FN56HV9AVYTX5N843` · parent [runtime](../runtime.md)

## acceptance

A testscript gets exit 2 with the next from rhiza request in an unregistered repo and in a repo under no root ($NOROOT, from testscripts-run-under-root), and runs rhiza check there successfully; the existing run start and request scripts pass under the Setup root.

## statement

Starting a run or queuing a request (rhiza run start, rhiza request) needs a registered project under a root and otherwise exits 2 with a next naming rhiza root init or rhiza project add; model commands (show, apply, check, render, land) keep working in any repo with .rhiza/, so CI needs no root.

## Links

### Out

- decided-by [m29-root-required](../../../decision/m29-root-required.md)

### In

- verifies [runs-need-root-verified](../../../evidence/runs-need-root-verified.md)
- verifies [runs-need-root-qa-1](../../../evidence/runs-need-root-qa-1.md)
- verifies [runs-need-root-verified-2](../../../evidence/runs-need-root-verified-2.md)
- verifies [runs-need-root-qa-2](../../../evidence/runs-need-root-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
