# plan-clean-round-skips-critic-step-qa-1

evidence · `01M77SRQ4N82JB6AQQ2YN1H1A6` · actor [claude-code-904cdffd](../actor/claude-code-904cdffd.md)

## result

PASS

## statement

QA at e520ec719 with the built binary (go build -trimpath into the session scratchpad): in a scratch Rhiza root and repo, with fake rex, claude and gh on PATH and a clean env, ran rhiza run start --foreground <record cs> --workflow plan (it runs plan-milestone). Run status showed review visit 1 (one finding) -> has-findings pass -> critic -> fresh -> check fail -> revise -> review visit 2 (clean {"findings":[]}) -> has-findings fail (exit 3: latest review round holds no finding, so it needs no critic) -> fresh -> check pass -> gate, with no critic visit 2, and prompts.log held exactly one critic prompt. Inside the live run, the fake reviewer ran rhiza run review-check --findings after each round: exit 0 after the round with a finding, exit 3 after the clean round. rhiza run verdict with {"verdicts":[]} on the clean round exited 2 ('at least one'). rhiza workflow describe plan and plan-milestone both show review on.pass -> has-findings. The run then failed at open-pr only because the hand-run fake gh lacks pr edit, which is outside this change.

## Links

### Out

- verifies [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
