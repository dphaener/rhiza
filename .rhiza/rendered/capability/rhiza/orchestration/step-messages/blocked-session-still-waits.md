# blocked-session-still-waits

capability · `01M65077HHPKVP9FVAECGBE1GJ` · parent [step-messages](../step-messages.md)

## acceptance

The existing session-message testscript still passes: with a question open no keys are sent, and the command waits for the prompt

## statement

A session blocked on an open question or menu is never typed into, and a message to it waits as it does today

## Links

### In

- verifies [blocked-session-still-waits-verified](../../../../evidence/blocked-session-still-waits-verified.md)
- verifies [blocked-session-still-waits-qa-1](../../../../evidence/blocked-session-still-waits-qa-1.md)
- verifies [pr-105-message-a-busy-session](../../../../evidence/pr-105-message-a-busy-session.md)

## History

- `01M61GW821FWSX05A9FC231N65` rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.
