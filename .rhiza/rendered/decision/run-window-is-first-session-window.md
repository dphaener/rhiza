# run-window-is-first-session-window

decision · `01M66JBKWS4MK7H2YBFSNDJRKP`

## reasoning

The built code already does this, and the run's first session is the first thing of the run with a pane, so opening the window there needs no window owned by daemonStart and keeps every session of one run, whatever its change set, side by side. Anchoring later sessions on the latest session with a pane, not the first, keeps them in the window the run opened last when an earlier one closed.

## statement

A run's panes live in one window per run, the one its first session with a pane opens (paneAnchor in internal/api/run_leaves.go). Each later session of the run splits its step anchor's tab when that has a record, else the tab of the run's latest session with a pane, which is the first session's while that window lives (see anchor-gone-splits-other-pane). Not the window rhiza run start (daemonStart) opens, and not one window per change set.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### Out

- answers [run-window-owner](../question/run-window-owner.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
