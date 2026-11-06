# m29-root-required

decision · `01M6RN37RJHQ058KK9BPVW4VY8`

## question

What happens to a repo that isn't under any root?

## quote

There isn't anything existing so root can be required from the start

## quote_where

milestone-29-interview, typed answer on a repo under no root

## statement

A root is required from the start: the daemon and runs work only on a registered project under a root, and no per-repo daemon mode is kept.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- decided-by [m29-root-registry](../intent/milestone-29/m29-root-registry.md)
- decided-by [m29-one-daemon](../intent/milestone-29/m29-one-daemon.md)
- decided-by [runs-need-root](../capability/rhiza/runtime/runs-need-root.md)
- decided-by [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)
- decided-by [testscripts-run-under-root](../capability/rhiza/runtime/testscripts-run-under-root.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
