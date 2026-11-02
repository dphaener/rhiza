# accept-policy

capability · `01M5SWTCJQP64M9ZDHGT7GS7PT` · parent [orchestration](../orchestration.md)

## statement

Acceptance policy decides which change sets are accepted without the owner: with accept auto, change sets an owner-requested plan run proposes are accepted at its gate when its review passes; with accept owner, they wait for the owner

## Links

### In

- verifies [accept-policy-qa-1](../../../evidence/accept-policy-qa-1.md)
- verifies [accept-policy-qa-3](../../../evidence/accept-policy-qa-3.md)
- verifies [pr-187-live-daemon-settings](../../../evidence/pr-187-live-daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
