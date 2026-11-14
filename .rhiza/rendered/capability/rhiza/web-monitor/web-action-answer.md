# web-action-answer

capability · `01M79Q8PZE56YE4MSC7P9GS5RT` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestWebAnswerFormsCarryCSRFAndAreAbsentWithoutToken|TestWebAnswerPostConfirms|TestWebApproveRefusesMissingCSRF|TestWebAnswerMergedInboxFormsPostIntoTheProject' shows the answer and approve forms carry the CSRF field and are absent when no token is set, answers a question and approves a sign-off through POST with the token and the CSRF field, shows the confirmation, and refuses an answer or approve POST missing the CSRF field; go test -trimpath ./internal/api -run 'TestWebAnswerAction|TestWebApproveAction' shows a POST through WebProjectHandler, and through WebRootHandler's merged inbox at the root, runs the real question.answer and inbox.approve or run.approve on the project's log.

## statement

From the web UI the owner answers an open question (question.answer) and approves a sign-off (inbox.approve or run.approve), and the page then says what happened ('Answered. The run continues.').

## Links

### Out

- depends-on [web-token](web-token.md)
- depends-on [web-csrf](web-csrf.md)
- decided-by [web-safe-actions](../../../decision/web-safe-actions.md)

### In

- verifies [web-action-answer-verified](../../../evidence/web-action-answer-verified.md)
- verifies [web-action-answer-qa-1](../../../evidence/web-action-answer-qa-1.md)
- verifies [web-action-answer-verified-2](../../../evidence/web-action-answer-verified-2.md)
- verifies [web-action-answer-qa-2](../../../evidence/web-action-answer-qa-2.md)
- verifies [pr-574-safe-actions-in-the-web-ui](../../../evidence/pr-574-safe-actions-in-the-web-ui.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
