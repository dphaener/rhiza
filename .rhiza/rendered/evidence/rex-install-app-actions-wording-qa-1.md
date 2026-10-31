# rex-install-app-actions-wording-qa-1

evidence · `01M686VGX30T6A0KCWW0R26N8E` · actor [claude-code-17358f03](../actor/claude-code-17358f03.md)

## result

PASS

## statement

At dccf0f0 ran the built binary with REX_CONFIG pointing at a scratch init.lua and a stub rex on PATH that only logs 'config reload': 'rhiza rex install' printed "Installed Rhiza's app actions and key mode in .../rhiza.lua (2 written, 0 unchanged); Rex reloaded.", a second run wrote 0 with 2 unchanged, --json reported the same, and the written rhiza.lua held no rhiza_launch; 'rhiza rex uninstall' printed "Removed Rhiza's app actions and key mode from Rex; Rex reloaded." and init.lua came back byte for byte (cmp). rhiza --help and describe --json show the new rex.install/rex.uninstall summaries, and describe --json holds 0 'quick-launch'. make check passed with RHIZA_ROOT/RHIZA_CHANGE/RHIZA_RUN unset.

## Links

### Out

- verifies [rex-install-app-actions-wording](../capability/rhiza/rex-install-command/rex-install-app-actions-wording.md)

## History

- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
