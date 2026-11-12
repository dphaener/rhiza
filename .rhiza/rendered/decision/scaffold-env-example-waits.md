# scaffold-env-example-waits

decision · `01M76C1M3TYKPZWHC55SZHQC97`

## question

The new-workflow scaffold and the docs are planned to show an example of setting environment variables in a workflow. That feature is built in a different change set, and the scaffold isn't required to wait for it, so it could ship examples of a field that doesn't exist yet. Should the scaffold wait for the environment-variable work, or should the environment example be dropped?

## statement

The language-neutral workflow scaffold waits for workflow environment variables and keeps its environment example

## wording

Wait and keep it: the scaffold change set depends on the environment-variable change set and keeps the example, picked by the owner

## Links

### Out

- answers [review-01m75w8yt6m1f8w05gajzsec0h-r3-f2](../question/review-01m75w8yt6m1f8w05gajzsec0h-r3-f2.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
