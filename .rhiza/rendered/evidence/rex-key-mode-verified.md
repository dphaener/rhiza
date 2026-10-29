# rex-key-mode-verified

evidence · `01M6482RM0675HPAG07KDCW0CD` · actor [claude-code-c61b6237](../actor/claude-code-c61b6237.md)

## result

PASS

## statement

At b411878: go test ./cmd/rhiza -run TestScripts/rex-install passes; the script now greps rhiza.lua for the launch key entering the mode (client.mode.enter, name = MODE), default keys for roadmap and model, settings keys = {}, and rex.bind(MODE .. "/" .. k, "rhiza_app_" .. app). Unit tests TestRexInstallBindsAKeyPerAppInTheKeyMode (one unique key per launcher app, each written in DEFAULT_KEYS), TestRexDefaultKeysAreMnemonicAndUnique and TestRexInstallKeepsKeyOverridesOnReinstall pass. A scratch rex install wrote 16 DEFAULT_KEYS entries (accept=a, change-show=c, ... why=w), all distinct, one for each of the 16 rhiza_app_* actions, in an exclusive mode "rhiza" with escape bound to exit. That matches decision rex-app-keys-default-overridable (default mnemonic keys, remapped through settings.keys, which a reinstall keeps) and native-rex-prefix-configurable (launch_key setting, default ctrl+a a). Rex docs confirm rex.mode exclusive and client.mode.enter { name, once }. gofmt, go vet and go test ./... pass. make check's only failure is the stale .rhiza/rendered from model writes, which this commit doesn't touch. Note: a settings.keys remap onto another app's default key (e.g. the comment's own example roadmap = "m" collides with model) binds both, and the winner depends on pairs order.

## Links

### Out

- verifies [rex-key-mode](../capability/rhiza/rex-key-mode.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
