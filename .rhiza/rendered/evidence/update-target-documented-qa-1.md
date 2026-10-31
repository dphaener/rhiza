# update-target-documented-qa-1

evidence · `01M6995JFG72C5NM8614PXCWWF` · actor [claude-code-083d55ec](../actor/claude-code-083d55ec.md)

## result

PASS

## statement

At 6eb01e4, make help lists 'update  Rebuild into $(PREFIX), then rewrite the Rex script and reload Rex' beside install. CLAUDE.md's Commands block has 'make update     # make install, then rewrite the Rex script and reload Rex'. TestMakeUpdateDryRun (it asserts make help names update) and TestClaudeMDCommandsNameMakeUpdate both PASSED with -count=1.

## Links

### Out

- verifies [update-target-documented](../capability/rhiza/rex-install-command/update-target-documented.md)

## History

- `01M686RWSEKXXBCNC7RZMT07HB` make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command
