# AgentJam Runbook — Operational Entry Point for AI Agents

You are operating in an AgentJam-managed workspace as **software-engineer**.
This file is your entry point. Follow the protocol below exactly.

## Enforcement Protocol (before you act)

1. Read the installed rules for your harness — they are domain-organized:
   ```
   skills/<category>/<skill>     — capability instructions
   policies/<category>/<policy>  — enforced rules (strict-block = must obey)
   agents/<id>                   — persona definitions
   ```
2. strict-block policies are non-negotiable. A response that violates one is
   a failed response.
3. Do not produce output that `agentjam scan` would flag.

## Scan & Fix Loop

```bash
agentjam scan              # policy-scan all source files; exit 1 on strict blocks
agentjam scan --offline    # skip package-registry freshness checks
agentjam eval <file>       # evaluate one file with line numbers
agentjam validate          # canonical resource integrity (errors must be 0)
agentjam run <workflow>    # structured multi-step workflows (e.g. bug-fixing)
agentjam init              # materialize canonical tree under .agentjam/ (editable source of truth)
agentjam uninstall         # remove installed files (ledger-tracked; --purge for full removal)
```

For every task:
1. **Scan first** — `agentjam scan` on the touched area before and after edits.
2. **Fix violations** by severity: strict-block → warning → info.
3. **Re-scan** until clean. `agentjam scan` exit code 0 is the gate.

## Tooling

```bash
agentjam mcp               # serve workspace tools (filesystem/git/terminal/run_workflow) over MCP stdio
agentjam preflight         # verify toolchain binaries for the detected stack
agentjam detect            # show detected harnesses and stacks
```

## Canonical Tree

If the project was initialized with `agentjam init`, the canonical
resource tree lives under `.agentjam/` (`agents/`, `skills/`, `policies/`, `tools/`, `workflows/` ...).
That tree is the source of truth — extend it rather than working around it,
and keep `agentjam validate` green. Projects without it still work:
scan falls back to the policies embedded in the binary.
