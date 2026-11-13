# pr-545-fix-rhiza-session-message

evidence · `01M7AENVW68W0MA5EKZ0NPPGMT` · actor [dphaener](../actor/dphaener.md)

## head

ac6dae2c80a5471463a8022b72a02a28532873d3

## pr

https://github.com/dphaener/rhiza/pull/545

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #545 (https://github.com/dphaener/rhiza/pull/545) passed CI at ac6dae2

## status

verified

## Links

### Out

- verifies [message-file-in-session-checkout](../capability/rhiza/orchestration/step-messages/message-file-in-session-checkout.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
