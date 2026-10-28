# AI usage and prompts

The prompts I used, in order, grouped by stage.

## 1. Planning

Claude Opus 5.5, Claude Fable 5.1 and GPT Sol 6.1, one session each.

> _The assignment brief, pasted unchanged._

## 2. Comparing and refining the plans

[Claude Code](https://claude.com/claude-code) on Claude Opus 5.5, with reviewer subagents on Claude Opus, Sonnet, Haiku and Fable.

> Compare these three plans against the assignment brief. Identify strengths, missing requirements, unnecessary complexity, and technical risks. Recommend one plan or a justified hybrid.

> Revisit the implementation estimates to account for AI-assisted development. Distinguish estimated coding effort from actual elapsed time, including execution, debugging, review, and verification.

> Produce the hybrid plan and have independent reviewers examine it. Verify uncertain technical claims against official documentation or executable examples, and address concrete findings before implementation.

## 3. Implementation

Claude Code on Claude Opus 5.5.

> Implement the approved plan step by step. Keep the domain logic independent of transport and presentation. Add meaningful tests alongside each layer, run the relevant checks before committing coherent changes, and document the API, setup, assumptions, and limitations.

## 4. Implementation review

Claude Code on Claude Opus 5.5, with reviewer subagents on Claude Opus, Sonnet, Haiku and Fable.

> Ask independent reviewers to evaluate the implementation using only the assignment brief and a fresh checkout. Do not provide the implementation plan or previous review conclusions. Follow the README, run the available checks, exercise the API, and inspect the interface. Distinguish reproduced defects from suspected issues and report the checks actually performed.
