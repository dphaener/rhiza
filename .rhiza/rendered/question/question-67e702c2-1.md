# question-67e702c2-1

question · `01M64JVMH1VSE24CBHWNGD9N8Q`

## session

67e702c2-a3f2-4cad-84dd-01a1892c6c98

## statement

The flake reproduces: under `go test ./internal/api -run TestGuideInstallConcurrent -count=300 -race -cpu 1,2,4,8` one run in 300 fails with "exists and was not written by rhiza guide install". From the code, pluginForeign checks plugin.json first and the other plugin files second, while another install writes plugin.json first and then the rest. So a racing install can see no manifest, then find SKILL.md written by the other install, and call the folder foreign. How should the fix go, without loosening the test?

## Links

### In

- answers [plugin-foreign-ignores-partial-installs](../decision/plugin-foreign-ignores-partial-installs.md)

## History

- `01M5ZNVHN2YNK6DDVYY8GK5CEQ` TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.
