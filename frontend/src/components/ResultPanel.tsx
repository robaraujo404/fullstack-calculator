import type { Ref } from 'react'

import { formatResult } from '../domain/formatResult'
import type { CalculatorState } from '../hooks/useCalculator'

// Friendlier copy for the errors a user can actually cause; anything else
// falls back to the message sent by the API.
const ERROR_MESSAGES: Record<string, string> = {
  DIVISION_BY_ZERO: 'Division by zero is undefined.',
  INVALID_OPERAND: 'The result is not a real number for these inputs.',
  RESULT_OUT_OF_RANGE: 'The result is too large to represent.',
}

interface ResultPanelProps {
  state: CalculatorState
  /** Expression for the current result, e.g. "10 ÷ 4". */
  expression: string
  ref?: Ref<HTMLDivElement>
}

export function ResultPanel({ state, expression, ref }: ResultPanelProps) {
  return (
    <div className="display" role="region" aria-label="Result" tabIndex={-1} ref={ref}>
      <p className="display-expression">{state.status === 'success' ? expression : ''}</p>
      {state.status === 'error' ? (
        <p className="display-error" role="alert">
          {ERROR_MESSAGES[state.code] ?? state.message}
        </p>
      ) : (
        <output className="display-result" aria-live="polite">
          {state.status === 'success' && formatResult(state.value)}
          {state.status === 'loading' && '…'}
        </output>
      )}
    </div>
  )
}
