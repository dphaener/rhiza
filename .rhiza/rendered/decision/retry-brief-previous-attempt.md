# retry-brief-previous-attempt

decision · `01M7QQ8DJSEWGBF0WGFGHCHG05`

## reasoning

The owner picked 'Right step + findings': the next attempt must see why the last one failed; brief.md is the file every step session already reads

## statement

Every agent step session started after a run's last retried entry gets a Previous attempt section in brief.md: the failed step, its outcome, its cause or output tail (as a hand-back reads them, run.CauseOf) and the findings of the latest round recorded before the retry, each with check, severity and text; the retried entry journals the failed step so the brief can name it. Spawned children's briefs are their caller's and are not changed

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-retry-conduct-plan](../source/m40-retry-conduct-plan.md)

### In

- decided-by [retry-brief-carries-failure](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-brief-carries-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
