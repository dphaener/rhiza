# review-01m6rmj6-r3-f3

question · `01M6RP4X9AKMA32RTMDSK9VRTQ`

## statement

Plan review round 3 left finding f3 (contradiction on usage-settings): The updated acceptance calls for 'A new usage-settings.txtar', but cmd/rhiza/testdata/script/usage-settings.txtar already exists. The check should say to update that script to cover machine scope and rhiza.toml. The critic confirmed it: cmd/rhiza/testdata/script/usage-settings.txtar already exists; the updated usage-settings acceptance says 'A new usage-settings.txtar', which contradicts the code. How should it be settled?

## Links

### Out

- blocks [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)

## History

- `01M6RP4XBJW9MQ8KZNRY8ZFAY8` Leave the questions of plan review round 3 of run 01M6RMJ6P7EG0VXT1752W0H8NY
