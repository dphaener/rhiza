# review-01m7ah6rayrzax6z6r63q0m2qk-r6-f4

question · `01M7BZ7EJEATB4NEHKEY1S89NV`

## statement

Plan review round 6 left finding f4 (contradiction on builtin-conductor-default): The conduct step's prompt has the conductor run rhiza ship rebase on a main-moved hand-back. But ship rebase (internal/api/ship_rebase.go) merges origin/main, runs the local gate and the pre-push check, and pushes the branch, recording the pushed head. The conduct topic (internal/guide/text/conduct.md: 'Merge or push: the ship run does both') forbids the conductor to push. The leaf also says the push step pushes the new head and never re-pushes an unchanged one, but after ship rebase the head is already pushed. The critic confirmed it: Verified: origin/main internal/api/ship_rebase.go (header :6-16, Pushed field :61) merges origin/main, runs the local gate and pushes the merged head. internal/guide/text/conduct.md:50 lists 'Merge or push: the ship run does both' among what the conductor must not do. The leaf has the conductor run ship rebase on a main-moved hand-back and also says the push step pushes the new head, but ship rebase has already pushed it. Contradiction, major stands. How should it be settled?

## Links

### Out

- blocks [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)

### In

- answers [conduct-main-moved-merge-no-push](../decision/conduct-main-moved-merge-no-push.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
