# plugin-foreign-race-interview

source · `01M64JXKM4VK65VCSYB5D6PDS4`

## quotes

```json
[]
```

## raw_sha256

dfb12a18c83e0f12b06d2e66c283a936894ca25089bf05a54daebdda394c55f1

## session

67e702c2-a3f2-4cad-84dd-01a1892c6c98

## summary

The concurrent first-install test fails about one run in 300 under -race with -cpu 1,2,4,8, with the error that the folder exists and was not written by rhiza guide install. The cause is in pluginForeign: it checks plugin.json first and the other plugin files second, while a racing install writes plugin.json first and then the rest, so a partial install can be called foreign. The question of how to fix it was put to the arbiter, since the run is autonomous and no owner answers. The arbiter decided that a folder with no plugin.json is foreign only when one of its other plugin files holds bytes this rhiza would not write. The install keeps writing plugin.json first and the test is not loosened.

## title

TestGuideInstallConcurrent flake: plan interview

## Links

### In

- derived-from [plugin-foreign-ignores-partial-installs](../decision/plugin-foreign-ignores-partial-installs.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
