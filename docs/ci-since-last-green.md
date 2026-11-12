# Skip the full CI round for model-only changes

Once `rhiza ci install` has given a repo the Rhiza check, the project's own test jobs still run in full on every push to a pull request. Most pushes to a pull request after the first don't need that. If nothing outside `.rhiza/` changed since the last commit of the pull request whose CI passed, the code and its tests are the ones that already passed, and only the model check has anything new to look at. A push that lands a plan, a render or a merge of the main branch that brought in only `.rhiza/` changes is the usual case.

This page is a recipe for a `changes` job that makes that decision. `rhiza ci install` does not write it, because it can't know your test jobs. You add it to your own CI.

## The rule

On a pull request, take the pull request's latest earlier commit whose CI passed: the **last green commit**. Then:

1. If there is no last green commit, run everything.
2. If the last green commit is not an ancestor of the head (after a force push or a rebase), run everything.
3. List the paths that differ between the last green commit and the head. If every one is under `.rhiza/`, run only the model check.
4. Otherwise run everything.

Anything that isn't a pull request (a push to the main branch, a manual run, a schedule) runs everything. The model check runs either way; "only the model check" means the other jobs are skipped.

Skipped is not failed, but it must not read as green on its own either. The job your branch protection requires has to pass only when the `changes` job succeeded and at least one job besides it ran.

## An example `changes` job

The example is for GitHub Actions. It names no language or test tool, and every step is a `git` or `gh` call, so it fits any project. Replace `test` with your own jobs.

```yaml
jobs:
  changes:
    runs-on: ubuntu-latest
    outputs:
      code: ${{ steps.diff.outputs.code }}
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    - name: Decide whether anything outside .rhiza/ changed since the last green commit
      id: diff
      env:
        GH_TOKEN: ${{ github.token }}
        EVENT: ${{ github.event_name }}
        PR: ${{ github.event.pull_request.number }}
        CHECK: ci-passed
      run: |
        full() { echo "code=true" >> "$GITHUB_OUTPUT"; exit 0; }
        [ "$EVENT" = pull_request ] || full
        head=$(git rev-parse HEAD)
        # The pull request's commits, newest first, the head excluded.
        for sha in $(gh api "repos/$GITHUB_REPOSITORY/pulls/$PR/commits" --paginate --jq '.[].sha' | tac | grep -v "^$head$"); do
          state=$(gh api "repos/$GITHUB_REPOSITORY/commits/$sha/check-runs" \
            --jq "[.check_runs[] | select(.name == \"$CHECK\")] | if length > 0 and all(.conclusion == \"success\") then \"success\" else \"\" end")
          if [ "$state" = success ]; then green=$sha; break; fi
        done
        [ -n "${green:-}" ] || full
        git merge-base --is-ancestor "$green" HEAD || full
        outside=$(git diff --name-only "$green" HEAD | grep -v '^\.rhiza/' || true)
        if [ -n "$outside" ]; then full; fi
        echo "code=false" >> "$GITHUB_OUTPUT"

  model:
    needs: changes
    if: needs.changes.outputs.code == 'false'
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - run: rhiza check

  test:
    needs: changes
    if: needs.changes.outputs.code == 'true'
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - run: ./run-your-tests   # your own build and test commands

  ci-passed:
    needs: [changes, model, test]
    if: always()
    runs-on: ubuntu-latest
    steps:
    - name: Every job succeeded or was skipped, and one besides changes ran
      env:
        RESULTS: ${{ join(needs.*.result, ' ') }}
        CHANGES: ${{ needs.changes.result }}
      run: |
        [ "$CHANGES" = success ] || exit 1
        ran=0
        for r in $RESULTS; do
          if [ "$r" = success ]; then ran=$((ran + 1)); elif [ "$r" != skipped ]; then exit 1; fi
        done
        [ "$ran" -ge 2 ] || exit 1
```

Make `ci-passed`, the aggregate job, the one required status. Its name is what `CHECK` looks for, so the last green commit means "its `ci-passed` passed".

Give the aggregate job a name no other job in the repo uses. GitHub names a check run by its job name alone, not by its workflow, and the workflow `rhiza ci install` writes already has a job named `check`. If the aggregate job were also named `check`, a commit whose tests failed but whose Rhiza check passed could be taken as green, and a later `.rhiza/`-only push would skip the tests that failed. The example also counts a commit green only when every check run of that name succeeded, so a clash fails safe: it runs everything.

## Things to know

- A merge of the main branch into the pull request counts as a change of everything the merge brought, because the diff runs from the last green commit. A merge that only brings `.rhiza/` changes from main stays model-only.
- Your own cheap paths (documentation, say) can join `.rhiza/` in the `grep -v` pattern. Keep that list to paths no test depends on.
- When the pull request has no green commit yet, the first push runs everything, which is the safe default.
