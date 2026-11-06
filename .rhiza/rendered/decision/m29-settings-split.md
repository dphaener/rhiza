# m29-settings-split

decision · `01M6RN37FV21SYDTNZKY4VP4AM`

## question

Where do the daemon's settings live once there is one root daemon?

## statement

Machine-wide daemon settings (maxRuns, usage budget and reserve, minFree, web host and port, selfUpdate, permissionMode, a global pause) live in the root config; per-project settings (autonomy per stage, project pause, per-project max, deploy check) stay in each repo's git config; rhiza daemon set writes to the scope the key belongs to.

## wording

Split machine/project. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-29-interview](../source/milestone-29-interview.md)

### In

- decided-by [m29-one-daemon](../intent/milestone-29/m29-one-daemon.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
