# testscripts-run-under-root-qa-2

evidence · `01M6S8XB8YF8TZP7QFBYTX0ADE` · actor [claude-code-9c5b9523](../actor/claude-code-9c5b9523.md)

## result

PASS

## statement

QA cycle 2: go test -trimpath ./cmd/rhiza -run 'TestScripts/^(root-init|project-add-path|project-add-url|project-ls-rm|runs-need-root|testscript-root|run-start|run-start-queue|daemon-settings|daemon-pause|daemon-event-hook|stage-autonomy)$' passes at HEAD 0099b258.

## Links

### Out

- verifies [testscripts-run-under-root](../capability/rhiza/runtime/testscripts-run-under-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
