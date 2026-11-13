While CI is unstable, or a red is not explained by your change,
do not push or rerun CI to see. Reproduce the
failing test locally, fix it, and push only after it and the
touched packages pass locally.
Reproduce flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host.
