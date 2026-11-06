# review-01m6rmj6-r3-f1

question · `01M6RP4X0H7YY4B9HFEGFYBBFY`

## statement

Plan review round 3 left finding f1 (contradiction on testscripts-run-under-root): The leaf registers each script's $WORK folder as a project, but the scripts it must keep passing unedited (run-start, run-start-queue, daemon-repo-session, daemon-event-hook, ...) do `mkdir repo; cd repo; git init` and work in $WORK/repo. $WORK is not a git checkout, and $WORK/repo sits two levels under the root, so project-add-path refuses both (non-git folder; nested path). The repos those scripts use would not be registered projects, so runs-need-root and root-daemon-lock would exit 2 there. The leaf has to register the repo folder the script actually uses, or let the scripts change. The critic confirmed it: run-start, run-start-queue, daemon-settings, daemon-pause, daemon-event-hook, daemon-repo-session and stage-autonomy all do mkdir repo; cd repo; git init and work in $WORK/repo. testscripts-run-under-root registers $WORK (not a git checkout, and the repo sits a level below it), which project-add-path says is refused (non-git folder, nested deeper). The repos these scripts use would not be registered projects, so runs-need-root/root-daemon-lock would exit 2 and the 'unedited' acceptance cannot hold. How should it be settled?

## Links

### In

- answers [m29-test-root-per-script](../decision/m29-test-root-per-script.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
