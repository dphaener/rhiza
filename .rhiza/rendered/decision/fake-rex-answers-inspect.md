# fake-rex-answers-inspect

decision · `01M7FT142VT43K432BXE64EJXH`

## reasoning

strace of bugfix-workflow-autonomous: 3,210 sleep 0.025 children of the fake rex. rhiza run start calls rex block inspect when REX_BLOCK is set (runProcess); these fakes have no case for it, so it falls into the ready-file loop meant for a started agent and times out after 1,600 polls (about 59 s a run start), and runProcess then ignores the failure. The wait is in the fake, not on rhiza's clock, so the scaled clock could never reach it. Answering inspect also exercises the window-id path runProcess was written for.

## statement

The scripts' fake rex answers rex block inspect at once with a window id, and waits for an agent's ready file only after it has started a pane's command (a call carrying --); any other call it does not know exits at once. The scripts keep every assertion they make today.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [fake-rex-no-dead-wait](../capability/fake-rex-no-dead-wait.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
