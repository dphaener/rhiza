# plugin-foreign-ignores-partial-installs

decision · `01M64JWR9ZY731P0ND1HK0AB3N`

## reasoning

The flake is a race between two guide.install runs of the same binary. Install A writes plugin.json first, then SKILL.md and the rest. Install B can stat the folder between those writes, see no plugin.json, find SKILL.md, and call the folder foreign. pluginForeign's comment already says a folder created by another install with no plugin.json yet should not be foreign, so the code intends this case and only misses it when a sibling file exists. Comparing the other files to pluginFiles(ctx.Version) keeps the real protection: a folder with a different SKILL.md or hooks file is still a conflict unless --force. Alternatives I weighed: writing plugin.json last (still racy, since the sibling-file check would still fire), and a lock file or flock around the install (new mechanism for a one-binary race). The fix stays inside guide_install.go and needs no change to the test.

## statement

pluginForeign treats a folder with no plugin.json as foreign only when one of its other plugin files holds bytes this rhiza would not write; a folder whose files all match this rhiza's own output is a racing install's partial write and is not foreign. The install keeps writing plugin.json first and does not loosen the test.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-67e702c2-1](../question/question-67e702c2-1.md)
- derived-from [plugin-foreign-race-interview](../source/plugin-foreign-race-interview.md)

### In

- decided-by [guide-install-concurrent-first](../capability/rhiza/guide/guide-install-concurrent-first.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
