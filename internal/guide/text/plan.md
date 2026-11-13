# Plan

A plan run writes a change set's leaves into the model: a planner session, with the owner
or autonomous, then a plan review. This topic is an index; the guidance lives in four role
topics. Load the one that is yours, as a skill or with `rhiza guide <topic>`:

- Planner: `/rhiza-planner` (`rhiza guide planner`). It reads `rhiza guide plan-recipes` and
  `rhiza guide plan-interview`.
- Plan reviewer: `/rhiza-plan-reviewer` (`rhiza guide plan-reviewer`).
- Plan critic: `/rhiza-plan-critic` (`rhiza guide plan-critic`).
- Plan reviser: `/rhiza-plan-reviser` (`rhiza guide plan-reviser`).

The planner runs in the plan step; the reviewer, critic and reviser run, in that order, in
the plan review (`rhiza guide plan-review`). Each role ends with `rhiza session done` as its
own command, as its topic says.
