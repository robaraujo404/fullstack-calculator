export type OperationId = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt' | 'percentage'

export interface Operation {
  id: OperationId
  /** Short text shown on the selector button. */
  symbol: string
  /** Accessible name of the selector button. */
  label: string
  arity: 1 | 2
  /** Renders the expression shown above the result. */
  expression: (a: string, b: string) => string
}

/** Single source for the selector, the accessible names and the API call. */
export const OPERATIONS: readonly Operation[] = [
  { id: 'add', symbol: '+', label: 'Add', arity: 2, expression: (a, b) => `${a} + ${b}` },
  { id: 'subtract', symbol: '−', label: 'Subtract', arity: 2, expression: (a, b) => `${a} − ${b}` },
  { id: 'multiply', symbol: '×', label: 'Multiply', arity: 2, expression: (a, b) => `${a} × ${b}` },
  { id: 'divide', symbol: '÷', label: 'Divide', arity: 2, expression: (a, b) => `${a} ÷ ${b}` },
  { id: 'power', symbol: 'xʸ', label: 'Power', arity: 2, expression: (a, b) => `${a} ^ ${b}` },
  { id: 'sqrt', symbol: '√x', label: 'Square root', arity: 1, expression: (a) => `√${a}` },
  {
    id: 'percentage',
    symbol: 'x% of y',
    label: 'Percentage: x% of y',
    arity: 2,
    expression: (a, b) => `${a}% of ${b}`,
  },
]

export function getOperation(id: OperationId): Operation {
  const op = OPERATIONS.find((o) => o.id === id)
  if (!op) throw new Error(`unknown operation: ${id}`)
  return op
}
