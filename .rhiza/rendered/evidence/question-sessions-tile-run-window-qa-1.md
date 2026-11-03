# question-sessions-tile-run-window-qa-1

evidence · `01M68KC8M2ERENW0VWGAMCPZFA` · actor [claude-code-881745c7](../actor/claude-code-881745c7.md)

## result

PASS

## statement

Same real-Rex autonomous plan run at 6b5628b. The planner asked two questions through the real PreToolUse hook. Each arbiter session opened with `rex split` of the planner's block, the second after the first arbiter's closed block failed as anchor and Rhiza split a live pane of the window instead. Neither arbiter opened a window, and the run succeeded with assume-1 and assume-2 landed. In a normal `ask` workflow run on real Rex, the advisor session for the owner question split the asking session's block in the run's one window, and rhiza inbox showed the recommendation.

## Links

### Out

- verifies question-sessions-tile-run-window (retired)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
