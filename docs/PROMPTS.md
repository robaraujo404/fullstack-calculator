# AI use and prompts

These are my prompts, in the sequence that I used them. Text in brackets gives context.

| Stage | Tools and models |
|---|---|
| Plans | Claude Opus 5.5, Claude Fable 5.1 and GPT Sol 6.1. Each model had a separate session. |
| All other stages | [Claude Code](https://claude.com/claude-code) on Claude Opus 5.5. The reviews used subagents on Claude Opus, Sonnet, Haiku and Fable. |

## 1. Plans

1. [The assignment brief, pasted unchanged into each of the three models.]

## 2. Plan comparison

2. "Run agents with different models to evaluate the plans: [paths to the three plans]. Score all three and see which one is best, or whether a hybrid version of them would be better."
3. "Too much scope is relative, I'm going to use AI to build it. So it can cover more hours." [The brief's 2–4 hour estimate assumes code written by hand. I kept the optional operations, the Dockerfile and CI, but nothing outside the brief: a calculation history was cut after all four reviewers flagged it as noise.]
4. "Write the hybrid and ask for another review, increasing the hours."
5. "What's missing for it to be a 9.5 or 10?" [The reviewers scored each plan out of 10.]
6. "Try to make those fixes and run the reviewers again."
7. "Apply them and run another round with only Sonnet and Opus." [After five rounds the plan scored about 8.9 (Sonnet 9.0, Opus 8.8). The late suggestions kept reversing each other, so I stopped reviewing the plan there.]

## 3. Implementation

8. "OK, go ahead." [Implement the final plan. It set the API contract, the numeric rules, the tests for each layer and the acceptance checks, and it ran as one commit per step: scaffold, calculator domain, REST API, API client and validation, UI, Makefile and coverage, README. The plan is not in the repository; the README's Design decisions section has its main decisions.]

## 4. Review

9. "Ask other models to review the implementation using only the assignment brief."
10. "Yes." [Fix the reproduced bugs and make the small improvements from the review, in one commit (3631c46). For example: a null body now gets 400 INVALID_REQUEST, POST /api gets a JSON 404 instead of an empty 307, percentage no longer overflows on representable results, and the UI rejects 1e-400 instead of sending 0. encoding/json's case-insensitive and repeated keys were kept as they are, pinned with tests and documented in the README.]
