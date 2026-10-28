# AI usage and prompts

This project was built with [Claude Code](https://claude.com/claude-code). The main agent ran on Claude Opus 5.5; plan reviews used independent subagents on Claude Opus, Sonnet, Haiku and Fable.

I wrote my prompts in Portuguese. Each one is quoted verbatim, followed by a one-line English summary. Status-only messages (e.g. asking whether a run had finished) are omitted. Entries are added as the work happens, not reconstructed afterwards.

## 1. Planning

Three implementation plans were drafted in separate sessions from the assignment email.

> _TODO: add the prompts used to draft the three initial plans._

## 2. Comparing and refining the plans

| # | Prompt (verbatim) | Summary |
|---|---|---|
| 1 | "rode agentes com diferentes modelos avaliando os planos: [3 paths] de nota pros 3 e veja qual o melhor ou se uma versao hybrida deles seria melhor" | Have agents on different models score the three plans and say whether a hybrid would be better. |
| 2 | "excesso de escopo é relativo, vou usar IA pra fazer. Então pode abranger mais horas" | Scope is relative since I'm using AI; the plan can cover more hours. |
| 3 | "escreva i hibrido e peça avaliacao novamente aumentando as horas" | Write the hybrid plan and have it reviewed again with a larger time budget. |
| 4 | "o que falta pra ser 9.5 ou 10?" | What is missing for a 9.5 or 10? |
| 5 | "tenta fazer essas correcoes e rodar os avaliadores novamente" | Apply those fixes and run the reviewers again. |
| 6 | "pq nao espera aumentar se tinha dois erros?" | Why not expect a higher score if two errors were fixed? |
| 7 | "roda" | Run the next review round. |
| 8 | "aplica e rode outra rodada só com sonneet e opus" | Apply the final fixes and run one more round with Sonnet and Opus only. |

Each review round gave four subagents the same rubric, written by the main agent, with the plans anonymized and shuffled. The rubric scored requirement coverage, proportionality (noise, unjustified pieces), technical design, testing, executability and evaluator impression, and asked reviewers to verify technical claims against official docs or by running code.

**Outcome.** Five rounds took the hybrid plan from 8.7 to 8.9 out of 10; the three original plans stayed between 6.0 and 7.6. Decisions that came out of the reviews:

- Keep all seven operations (the brief lists the advanced ones as optional), a single Dockerfile and CI; leave out anything the brief doesn't ask for, such as a calculation history.
- Every response under `/api/` is JSON, including 404, 405 and 500. The handler checks the method itself, because Go's `ServeMux` returns plain-text 405s and stops returning them at all once a `/` pattern serves static files.
- 400 means a malformed request; 422 means a well-formed request that is mathematically invalid.
- `percentage` is `a * b / 100`, which keeps common cases exact (`7% of 100 = 7`).
- The UI rounds non-integers to 15 significant digits for display and never alters integers.
- The Docker image is validated by a CI smoke test, since Docker isn't installed locally.

## 3. Implementation

The agent implemented the plan step by step: backend domain, HTTP layer, frontend logic, UI, Docker and CI, then docs. Each step was committed only after its tests, linters and type checks passed. It also launched the real app and drove it in headless Chrome to check the UI on desktop and at 375 px, and ran every `curl` example in the README against the server before including it.

| # | Prompt (verbatim) | Summary |
|---|---|---|
| 9 | "ok, segue" | Go ahead with the implementation. |
| 10 | "dentro da pasta ~/apps/pessoal tem o ssh pra minha conta pessoal robaraujo404" + repository name "fullstack-calculator" | Publish under my personal GitHub account (its SSH key is already configured) as `fullstack-calculator`. |
