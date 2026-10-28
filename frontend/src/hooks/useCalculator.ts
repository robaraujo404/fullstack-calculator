import { useCallback, useState } from 'react'

import { ApiError, calculate, type CalculateFn, type Operands } from '../api/calculatorApi'
import type { OperationId } from '../domain/operations'

export type CalculatorState =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; value: number }
  | { status: 'error'; code: string; message: string }

/** The form is disabled while a request runs, so it must never hang. */
export const REQUEST_TIMEOUT_MS = 10_000

/**
 * Runs calculations against the API and tracks the request state.
 * The client is injectable so tests don't need to mock fetch.
 */
export function useCalculator(client: CalculateFn = calculate) {
  const [state, setState] = useState<CalculatorState>({ status: 'idle' })

  const run = useCallback(
    async (op: OperationId, operands: Operands) => {
      setState({ status: 'loading' })
      try {
        const value = await client(op, operands, { signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS) })
        setState({ status: 'success', value })
      } catch (err) {
        const { code, message } =
          err instanceof ApiError ? err : { code: 'UNEXPECTED_ERROR', message: 'Something went wrong.' }
        setState({ status: 'error', code, message })
      }
    },
    [client],
  )

  const reset = useCallback(() => setState({ status: 'idle' }), [])

  return { state, run, reset }
}
