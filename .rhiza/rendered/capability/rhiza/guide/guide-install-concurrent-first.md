# guide-install-concurrent-first

capability · `01M64JXSEJKRWW6GVN295ZX1XY` · parent [guide](../guide.md)

## acceptance

go test ./internal/api -run 'TestGuideInstallConcurrent$' -count=300 -race -cpu 1,2,4,8 passes, with TestGuideInstallConcurrent unchanged, and make check passes

## statement

Concurrent first installs of the agent guide all succeed: a folder another install is still writing is not called foreign

## Links

### Out

- decided-by [plugin-foreign-ignores-partial-installs](../../../decision/plugin-foreign-ignores-partial-installs.md)

### In

- verifies [guide-install-concurrent-first-verified](../../../evidence/guide-install-concurrent-first-verified.md)
- verifies [guide-install-concurrent-first-qa-1](../../../evidence/guide-install-concurrent-first-qa-1.md)
- verifies [pr-83-testguideinstallconcurrent](../../../evidence/pr-83-testguideinstallconcurrent.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
