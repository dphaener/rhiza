# web-action-answer-qa-1

evidence · `01M7BK2AKZ8W4BA00ASQ50VG79` · actor [claude-code-a17d514a](../actor/claude-code-a17d514a.md)

## result

FAIL

## statement

QA cycle 1 with the built binary as a scratch root daemon: a real run stopped at a sign-off showed an Approve form in the root merged inbox posting to /p/proj/actions/approve; approving there journalled one signoff and queued the continuation ('Approved. The run continues.'), a second approve and an unknown run got 'Nothing to do', an empty run 'Nothing to approve', no csrf 403. A real rhiza run gate stop on two open questions showed answer forms (statements HTML-escaped, handles only in hidden fields); answering wrote the answering Decision in the run's worktree and a blank answer was refused. But the run never continued after both were answered: it stayed stopped at its gate, the questions and their forms stayed in both inboxes, and a repeat answer wrote a second Decision (finding f1).

## Links

### Out

- verifies [web-action-answer](../capability/rhiza/web-monitor/web-action-answer.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
