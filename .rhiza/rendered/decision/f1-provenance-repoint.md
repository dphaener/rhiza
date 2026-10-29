# f1-provenance-repoint

decision · `01M64587C3Z5R6QKX6ECJTWT7C`

## reasoning

I read the question, the decision rex-app-own-action, and both Sources. The quote the decision carries (the forge-cli line) is verbatim in native-rex-launch-owner-words, whose quotes are by owner and where the orchestrator session says it was captured. native-rex-plan-interview holds quotes [] and its summary says no owner words were captured this run, so the derived-from link points at a Source with no owner words for this quote. The defect is only the link target, not the words, so the fix is a re-point. Arbiter writes are limited to this Decision and its answers link, so the unlink of the old derived-from and the new derived-from to native-rex-launch-owner-words must be made by the builder or planner in change set 01M643X7 (the plan-interview Source may still be linked as the plan run's source).

## statement

Settle finding f1 by re-pointing rex-app-own-action's derived-from link from native-rex-plan-interview to native-rex-launch-owner-words, the Source that holds the owner's quote; the plan-interview Source stays as the record of the autonomous plan run.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m644m1-r3-f1](../question/review-01m644m1-r3-f1.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
