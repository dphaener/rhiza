# m22-leaves-split

decision · `01M6D4SCEZY3S2VD13P81PQDFS`

## statement

A run builds only leaf capabilities, so plan-requests and health-command, which were given children, would never be built: request-queue (in its own change set after plan requests) and health-report (in the health change set) carry their work as leaves

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m22-run-from-rex-ask](../source/m22-run-from-rex-ask.md)

## History

- `01M6D4SA06J84GXHJH8EX3W597` Fix milestone 22's leaves: the request queue and health report were grouping nodes with children, so runs would never build them
