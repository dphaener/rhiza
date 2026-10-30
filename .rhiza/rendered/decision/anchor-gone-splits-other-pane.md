# anchor-gone-splits-other-pane

decision · `01M66JBMPAHK74J6ZZMZT5KVZ8`

## reasoning

The question's recommended answer kept the run in its one window when a single pane closes. Review cycle 2 (qa f1, f2) showed the run's window goes away when its last pane closes, at a step's end or when agents exit, so failing there broke every multi-step run; opening a window then is a recoverable placement.

## statement

When a session's anchor pane is gone, the new session splits another live pane of the anchor's window (livePane, from rex block ls). Only when Rex no longer has that window, its last pane closed, or the window has no other pane, does the session open a new window; a session with no anchor of its own anchors on the run's latest session with a pane, so the run's later sessions join that new window.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### Out

- answers [anchor-gone-policy](../question/anchor-gone-policy.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
