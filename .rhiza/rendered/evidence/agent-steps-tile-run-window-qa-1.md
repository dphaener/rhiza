# agent-steps-tile-run-window-qa-1

evidence · `01M68KC7M4BEN4TDK491C6VRS2` · actor [claude-code-881745c7](../actor/claude-code-881745c7.md)

## result

PASS

## statement

At 6b5628b I built rhiza into a scratch dir and ran `rhiza run start <record> --workflow plan --autonomous` in a scratch repo against real Rex, in a dedicated Rex session. A caller-side rex shim forwarded every call to the real rex and only pointed new panes at a fake claude that stays alive after `rhiza session done`, as Claude Code does. Real Rex got exactly one `window new` (interview). The review session split the interview's live block, the critic split the review's block, and all five session records share one window. Panes left open by `until held` stayed alive. When the interview window had been emptied (an agent that exits on its own), the next step correctly fell back to `window new`. The named Go tests and testscripts also pass.

## Links

### Out

- verifies agent-steps-tile-run-window (retired)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
