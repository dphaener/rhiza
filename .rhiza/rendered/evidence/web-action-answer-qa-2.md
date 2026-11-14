# web-action-answer-qa-2

evidence · `01M7BQCMR3FA1C96SXAJ8KRKY9` · actor [claude-code-ccbe6c6d](../actor/claude-code-ccbe6c6d.md)

## result

PASS

## statement

QA cycle 2 with the branch built to a scratch path and run as a scratch root daemon (web token set): a real command-only plan run stopped at rhiza run gate on two open questions. The root merged inbox showed both with answer forms posting to /p/proj/actions/answer (csrf hidden field, statement HTML-escaped). A cookie+csrf POST answering the first got 200 'Answered. The run continues once its other questions are answered.' and wrote one Decision in the run's worktree; repeating it got 400 'Nothing to do' and wrote no second Decision; no csrf 403; blank text 400. Answering the last got 'Answered. The run continues.', the daemon resumed the run and its next step ran: the run succeeded and the inbox emptied. A second run with one question answered through rhiza question answer --run and the other through the web also resumed and succeeded; a web repeat of the CLI-answered one and an unknown question were refused. Answering an ended run 409, an unknown run 404, a missing question 400, no auth 401. Approve on a non-stopped or unknown run 'Nothing to do', empty run 'Nothing to approve'. Scoped go test ./internal/api -run 'TestWebAnswerAction|TestWebApproveAction|TestWebPauseAction|TestWebTokenSetting' and ./internal/web -run TestWeb pass.

## Links

### Out

- verifies [web-action-answer](../capability/rhiza/web-monitor/web-action-answer.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
