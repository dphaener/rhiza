# roadmap-app-new-work-qa-1

evidence · `01M6A8S2X693P9EJSSRYACFAXP` · actor [claude-code-b0fb6bcc](../actor/claude-code-b0fb6bcc.md)

## result

PASS

## statement

QA at cafec91c: ran the scratch build's `rhiza tui app roadmap` in a 200x40 pty with a logging fake rex on PATH. m shows 'Intent of the new milestone: _', typed Go and enter show 'Start an interactive plan run of a new milestone "Go" in a new Rex tab? y/n', y makes one call `rex window new -s session:qa new-milestone -c <repo> --json -- <rhiza> run start --workflow plan --new Go --milestone` with no change set argument, and the app stays open saying 'Opened Rex tab running the plan run of a new "Go"'. n gives new-change-set and no --milestone; backspace edits the text (Shipx then backspace gave --new Ship). Empty enter, esc at the intent prompt, and n or esc at y/n made no call. Works on an empty roadmap. Seen: the tab's run start carries no --permission-mode, so with the built-in plan workflow it exits 2 in the tab ('agent step with no permission mode') before opening a draft. The existing p key has the same gap. Raised as a discovered-from draft change set.

## Links

### Out

- verifies roadmap-app-new-work (retired)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
