# question-session-panes-close-qa-1

evidence · `01M68KC9PAB6ERTJX71QJW33TP` · actor [claude-code-881745c7](../actor/claude-code-881745c7.md)

## result

PASS

## statement

On real Rex at 6b5628b, with fake agents that keep running after `session done`: each arbiter's journal had `arbitrated` (08:57:41.99Z), then rhiza issued `rex block close` of the arbiter's block, then sent the answer to the asker (`answered` 08:57:42.09Z). The arbiter block vanished from `rex block ls` within a second, killing its still-running agent, while the interview, review and critic panes stayed open. For the advisor, `recommended` was journalled at 09:01:22.018Z and `rex block close` of the advisor's block followed at once. The pane and its process were gone within a second, while the asker's pane stayed open. A first advisor run whose pane PATH lacked rex left the pane open (no close), confirming the close comes from run recommend.

## Links

### Out

- verifies [question-session-panes-close](../capability/rhiza/runtime/rex-sessions/question-session-panes-close.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
