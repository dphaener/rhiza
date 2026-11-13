# fixer-lint-section-tighten-keep-rules

decision · `01M7B892BS9CV8FADZYGD0QYAF`

## reasoning

Attempt 1's verifier found two problems: f1, existing fixer rules rewritten and dropped, and f2, the Review:/session-done sentence missing. Attempts 2 and 3 made no commit, because the related question fixer-lint-section-over-topic-cap showed that the restored topic costs 994 of 1000 tokens, so an unchanged topic leaves no room for the section. Raising TopicCap (option a) would loosen a limit on every guide topic, which goes beyond this change set's intent (model lints in review). A new lint-fix topic (option b) would add files and a pointer outside the leaf's named files, and it would still need trimming. Shortening wording while keeping every rule (option c) meets what the leaf is for: a fixer that keeps all its existing duties. It touches only the two files the leaf names. It also fixes both of attempt 1's findings directly, since f1 was about rules being lost, not about words changing. Attempt 1's diff shows the existing text can be shortened by roughly the ~120 tokens the section needs.

## statement

Build @fixer-lint-hits-guide inside the existing TopicCap of 1000 by tightening the fixer topic's wording, not by raising the cap or adding a topic. Read "the rest of the fixer topic is unchanged" as: every existing rule is kept, though its wording may be shortened. Attempt 1 (453848a9d) dropped rules, and the next attempt restores each of them in short form: read the code at each finding's place before changing it; fix it matching the surrounding code; commit only once the finding's commands and the touched packages pass locally; a finding that could not be fixed is listed as unfixed in the commit; and the closing line that names the unfixed findings and why before `rhiza session done`. The "A lint-fix fixer" section also says it ends like any fixer, with one commit trailed `Review:` and then `rhiza session done` (verifier finding f2). fixer_test.go checks for those restored rules as well as the lint-fix phrases, and the topic stays under TopicCap. The builder adds a fix commit on top of 453848a9d and does not rewrite it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m7b0sj-1](../question/leaf-failed-01m7b0sj-1.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
