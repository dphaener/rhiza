# m29-update-usage-settings-script

decision · `01M6RYCYBEGVT9KQV14WX4Y7DT`

## statement

The usage-settings acceptance updates the existing cmd/rhiza/testdata/script/usage-settings.txtar instead of adding a new one: it shows the three keys with defaults and machine scope, sets each and finds it in the root's rhiza.toml rather than git config, still rejects a reserve outside 0-100 and a non-number limit with exit 2, and regenerates describe-golden.txtar.

## wording

Update the existing usage-settings.txtar to cover the settings moving to the machine-wide rhiza.toml; the agent's option, picked by the owner

## Links

### Out

- answers [review-01m6rmj6-r3-f3](../question/review-01m6rmj6-r3-f3.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
