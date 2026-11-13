# tiered-verifier-plan

source · `01M79CYASZDCHSQBHG21G3NF5C`

## quotes

```json
[]
```

## raw_sha256

5a8be67fa1791511e838797159acd119d78ee0e0334a4e6dd231f83de0ff1ae6

## session

03f1e4cd-1fae-4504-9df6-01a1d2cbb3c5

## summary

Autonomous plan run for change set Tiered verifier (milestone 35, ship-change median at most 40 minutes). No owner took part. The planner read the code: every built-in verifier (internal/workflow/builtin/build.toml, review.toml) runs models = ["opus"], the model comes from workflow.Role.Model at launch in internal/run/leaves.go startLeaf, and the leaf-session journal entry already records the model. It raised one question in four parts: what makes a leaf risky, which model verifies a low-risk leaf, where the configuration lives, and which verifier sessions are covered. The arbiter answered with Decision tiered-verifier-rules: risk is decided deterministically at launch from the paths the commit under verification changes, matched against the role's risky_paths globs, or from a retry (attempt past 1, a fresh round, or a review cycle past 1); an unreadable diff counts as risky. Low-risk leaves verify on sonnet and risky ones stay on opus. models stays the strong ladder, with new low_risk_models and risky_paths fields on the verifier role, and the built-in build and review verifiers get defaults. A role without the new fields behaves as today. Planned as three leaves: the role fields and their description, the conductor's choice at launch, and the built-in defaults.

## title

Tiered verifier plan, autonomous

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
