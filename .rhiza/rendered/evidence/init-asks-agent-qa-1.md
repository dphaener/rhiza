# init-asks-agent-qa-1

evidence · `01M7EAC1WP55RWVPES3H2VR3AK` · actor [claude-code-0cae456c](../actor/claude-code-0cae456c.md)

## result

PASS

## statement

QA with a fresh scratch build: rhiza init --agent claude-code in a git repo writes rhiza.project.agent and says so; --agent bogus exits 2 with the registered names and creates no .rhiza; stdin /dev/null and --json take claude-code without asking and the response has agent and saved; outside git it writes no setting and prints the rhiza daemon set command. Through a pty (script): an empty line takes the default, bogus is refused and asked again, then claude-code is stored; on a pty with --json nothing is asked; re-init of an existing repo asks nothing and fails as a conflict (exit 3). TestScripts/(init|root-init|describe-golden) pass.

## Links

### Out

- verifies [init-asks-agent](../capability/rhiza/model/repo-init/init-asks-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
