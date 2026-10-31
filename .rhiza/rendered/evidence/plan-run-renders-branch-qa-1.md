# plan-run-renders-branch-qa-1

evidence · `01M6958033GAXY78KMVKY4N87K` · actor [claude-code-418647f6](../actor/claude-code-418647f6.md)

## result

PASS

## statement

QA cycle 1 at a3dd859: built rhiza with go build to a scratch path and ran it by hand in a scratch repo with a bare origin.git, a fake rex/claude (roles doing real rhiza writes) and a fake gh (no network). rhiza run start <record> --workflow plan succeeded through gate -> open-pr. The branch's last commit was 'Render the plan' with Change: <cs> and Actor: trailers, holding every .rhiza/log entry the run's agents wrote plus .rhiza/rendered; git status --porcelain --untracked-files=all was empty and rhiza check printed ok with no manual render.

## Links

### Out

- verifies [plan-run-renders-branch](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-renders-branch.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
