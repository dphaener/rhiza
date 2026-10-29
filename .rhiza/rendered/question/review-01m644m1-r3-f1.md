# review-01m644m1-r3-f1

question · `01M6457M3EA8MGH73R0AW7VNTJ`

## statement

Plan review round 3 left finding f1 (missing-provenance on rex-app-own-action): The decision's quote ("the forge-cli opens the forge quick actions entirely differently, with a lua script, natively in Rex.") is in native-rex-launch-owner-words (01M643X6FYHZDDRBEPMQTZ3FDG), but its only derived-from Source is native-rex-plan-interview (01M644QHPB5W2EYTATMCZFJ82T), which holds quotes: [] and says no owner words were captured this run. The quote is not traced to the Source that holds it, so a reader following derived-from reaches a Source with no owner words. The critic confirmed it: Decision rex-app-own-action (01M644QVPPBD4ZKTT1YP9SMGAS) carries a quote that is verbatim in native-rex-launch-owner-words (01M643X6FYHZDDRBEPMQTZ3FDG), so the words are the owner's and are not invented. But its only derived-from link goes to native-rex-plan-interview (01M644QHPB5W2EYTATMCZFJ82T), whose quotes are [] and whose summary says no owner words were captured this run. A reader following derived-from reaches no owner words, so the recorded provenance is wrong. The kind is a stretch (a derived-from Source exists); the defect is the derived-from pointing at the wrong Source. Fix: re-point derived-from to native-rex-launch-owner-words. How should it be settled?

## Links

### In

- answers [f1-provenance-repoint](../decision/f1-provenance-repoint.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
