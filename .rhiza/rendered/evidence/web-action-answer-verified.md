# web-action-answer-verified

evidence · `01M7BHQQTYDWNFGW0RCHCT22HP` · actor [claude-code-05bcd10d](../actor/claude-code-05bcd10d.md)

## result

PASS

## statement

At a1857908f re-ran every Checks: command, all exit 0: go test -trimpath ./internal/web -run 'TestWebAnswer|TestWebApprove' (4 tests PASS: forms carry the CSRF field and are absent with no token/wrong cookie/no token set; answer and approve POSTs with token+CSRF confirm 'Answered. The run continues.'/'Approved. The run continues.'; answer and approve POSTs missing or with a wrong CSRF field get 403 and run nothing; merged inbox forms post into /p/<project>); go test -trimpath ./internal/api -run 'TestWebAnswerAction|TestWebApproveAction' (PASS: POST through WebProjectHandler runs the real question.answer, writing the answering Decision, and the real run.approve journals one sign-off, a second approve refused without a run ID; through WebRootHandler the merged inbox at the root shows an approve form posting to /p/repo/actions/approve and the POST runs run.approve; no token gives 401); go test -trimpath ./internal/web; go test -trimpath ./internal/api -run 'Web|Inbox|Describe|Docs|Signoff|ReadOnly'; make self-check; rhiza ship generate; rhiza render. Read the diff: changes are limited to the web actions (internal/web actions, inbox/merged inbox forms, templates, css, api WebActions and handler wiring) plus the leaf's log and rendered actor; no stubs or TODOs; tests use fakes and temp stores only.

## Links

### Out

- verifies [web-action-answer](../capability/rhiza/web-monitor/web-action-answer.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
