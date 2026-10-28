import { useEffect, useRef, useState, type FormEvent } from 'react'

import type { CalculateFn } from '../api/calculatorApi'
import { getOperation, OPERATIONS, type OperationId } from '../domain/operations'
import { parseOperand } from '../domain/parseOperand'
import { useCalculator } from '../hooks/useCalculator'
import { OperandInput } from './OperandInput'
import { ResultPanel } from './ResultPanel'

interface FieldErrors {
  a?: string
  b?: string
}

interface CalculatorProps {
  /** Overrides the API client; used by tests. */
  client?: CalculateFn
}

export function Calculator({ client }: CalculatorProps) {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [operationId, setOperationId] = useState<OperationId>('add')
  const [errors, setErrors] = useState<FieldErrors>({})
  const { state, run, reset } = useCalculator(client)
  const resultRef = useRef<HTMLDivElement>(null)

  const operation = getOperation(operationId)
  const needsB = operation.arity === 2
  const loading = state.status === 'loading'

  // The fieldset is disabled while the request runs, which drops keyboard
  // focus. Move it to the outcome so keyboard and screen reader users land there.
  useEffect(() => {
    if (state.status === 'success' || state.status === 'error') {
      resultRef.current?.focus()
    }
  }, [state.status])

  // Any edit invalidates the shown result, so it never sits next to inputs
  // that did not produce it.
  function changeA(value: string) {
    setA(value)
    setErrors((prev) => ({ ...prev, a: undefined }))
    reset()
  }

  function changeB(value: string) {
    setB(value)
    setErrors((prev) => ({ ...prev, b: undefined }))
    reset()
  }

  function changeOperation(id: OperationId) {
    setOperationId(id)
    setErrors((prev) => ({ ...prev, b: undefined }))
    reset()
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const parsedA = parseOperand(a)
    const parsedB = needsB ? parseOperand(b) : null
    setErrors({
      a: parsedA.ok ? undefined : parsedA.error,
      b: parsedB?.ok === false ? parsedB.error : undefined,
    })
    if (!parsedA.ok || parsedB?.ok === false) {
      return
    }

    void run(operation.id, parsedB?.ok ? { a: parsedA.value, b: parsedB.value } : { a: parsedA.value })
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} noValidate aria-busy={loading}>
      <ResultPanel state={state} expression={operation.expression(a.trim(), b.trim())} ref={resultRef} />

      <fieldset className="controls" disabled={loading}>
        <legend className="visually-hidden">Calculation</legend>

        <OperandInput
          id="operand-a"
          label={needsB ? 'First number' : 'Number'}
          value={a}
          error={errors.a}
          onChange={changeA}
        />

        <fieldset className="operations">
          <legend className="visually-hidden">Operation</legend>
          {OPERATIONS.map((op) => (
            <label key={op.id} className="operation">
              <input
                type="radio"
                name="operation"
                value={op.id}
                checked={op.id === operationId}
                onChange={() => changeOperation(op.id)}
                aria-label={op.label}
              />
              <span aria-hidden="true">{op.symbol}</span>
            </label>
          ))}
        </fieldset>

        {needsB && (
          <OperandInput id="operand-b" label="Second number" value={b} error={errors.b} onChange={changeB} />
        )}

        <button type="submit" className="equals" aria-label="Calculate">
          =
        </button>
      </fieldset>
    </form>
  )
}
