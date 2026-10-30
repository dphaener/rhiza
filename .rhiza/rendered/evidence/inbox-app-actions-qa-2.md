# inbox-app-actions-qa-2

evidence · `01M6763X297R6D37JA7NZ2VJET` · actor [claude-code-ed59bd32](../actor/claude-code-ed59bd32.md)

## result

PASS

## statement

Built at 602399e (includes fix 92739bb) with go build outside PATH and drove rhiza tui app inbox in a real Rex window with rex send and rex capture, over a scratch repo (RHIZA_ROOT/CHANGE/RUN/ACTOR unset) holding a question run and real review runs started with rhiza run start, whose signoff passes to a 20s wait step; a fake rex on PATH logs its argv and runs a window new command in the background, as a tab would. Pin: a, typed 'use pg', then removed the question's journal; the next tick showed 'that item no longer waits', left answer mode and sent no session.message. y on a sign-off ran rex window new -s <session> 01M675XE1 -c <repo> --json -- <binary> run approve <run> and showed 'approve opened in a new tab' at once; j moved the cursor within 0.5s while the tab's run sat in the wait step, so the app no longer freezes; the run journalled approved, finished after the wait, and the item left the list on the next tick. c with 'add tests: "x" & $HOME 100%' ran sh -c 'printf %s "$2" | "$1" run request-changes - "$3"' with the findings JSON as an argument; the run journalled outcome changes with that reason byte for byte and went back to review. R on a failed run ran run resume <run> --retry in a new tab; x then n said not abandoned, x then y abandoned it and it left the list; t said the advisor has no recommendation yet; an answer went to session S1 (fails: the fixture session has no record). Known minor: the owner's finding is filed under check quality, so on a workflow without one the tab fails with 'check "quality" is not a check of step review' while the app says the tab opened.

## Links

### Out

- verifies [inbox-app-actions](../capability/rhiza/inbox-app-actions.md)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
