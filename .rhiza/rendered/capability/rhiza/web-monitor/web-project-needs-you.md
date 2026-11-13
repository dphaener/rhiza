# web-project-needs-you

capability · `01M7B25DXMVDVRKGHGYNPVZGQP` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web ./internal/api -run 'TestProjectNeedsYou' passes: for a project fixture with one open question, one sign-off, one failed run, one run in an agent step whose last journal entry is 21 minutes old and one hold, the ProjectView's needs-you count is 5; a run whose last journal entry is 5 minutes old adds nothing, and a project with nothing waiting has a count of 0 (failure case: a run with no journal entry for 19 minutes counted as stuck fails).

## statement

Each project's needs-you count is a new ProjectView field (internal/web/projects.go), filled by its producer in internal/api/web_root.go from that project's model and runs, and defined here: the number of open questions, change sets waiting for the owner's sign-off, failed runs, stuck runs (a run in an agent step with no journal entry for 20 minutes) and holds. One function computes it and lists the items, so web-home-needs-you (change set 01M79NM3RJZFSZE1XHGHCRZSRM, which depends on this change set and builds after it) can count and list Home's 'things need you' with the same function, and the header, the picker and Home's title show the same number. A project's live state is Running plus Waiting above zero.

## Links

### Out

- decided-by [web-project-picker](../../../decision/web-project-picker.md)

### In

- depends-on [web-picker-sheet](web-picker-sheet.md)
- depends-on [web-frame-header](web-frame-header.md)
- blocks [review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f1](../../../question/review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f1.md)
- blocks [review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f2](../../../question/review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f2.md)
- verifies [web-project-needs-you-verified](../../../evidence/web-project-needs-you-verified.md)
- verifies [web-project-needs-you-qa-1](../../../evidence/web-project-needs-you-qa-1.md)
- verifies [web-project-needs-you-qa-2](../../../evidence/web-project-needs-you-qa-2.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
