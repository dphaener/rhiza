# pr-501-command-budget

evidence · `01M78JQQQHV4KEZBBY3KGHNTB2` · actor [dphaener](../actor/dphaener.md)

## head

ecde694504e29e62917bc2f639fb369e499fc239

## pr

https://github.com/dphaener/rhiza/pull/501

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #501 (https://github.com/dphaener/rhiza/pull/501) passed CI at ecde694

## status

verified

## Links

### Out

- verifies [command-budget-guide](../capability/rhiza/guide/command-budget-guide.md)

## History

- `01M76XQX53GGNWVCQHXVV510NH` Builder, verifier, check and fixer guidance says: never sleep or poll, no single command longer than about four minutes, and a long command runs once in the background and is waited for once (afk execute.js rule). Evidence: verifier session p90 33.7m against a 3.7m median, builder p90 25-36m against about 4m (prompt-gap rec 6).
