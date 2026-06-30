# Grounding (Foundry IQ)

Retrieved **3** source(s) to ground this run via **foundry**.

## Agentic retrieval — planned sub-queries
- enable OAuth support in authentication service for external tokens

## Sources
- **[S1] S-4 Secrets handling** — `security-checklist.md`
  - Secrets (API keys, tokens) come from environment variables or a secret store — never hardcoded, committed, or returned to the browser.
- **[S2] S-2 Authentication errors** — `security-checklist.md`
  - Failed logins must return a generic `401` ("invalid email or password") that does not reveal whether the email exists. Avoid user-enumeration leaks.
- **[S3] C-5 Dependencies** — `coding-standards.md`
  - Add dependencies deliberately and pin sensible minimums. Prefer the standard library and already-present packages.