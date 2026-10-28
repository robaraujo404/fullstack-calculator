export type OperationId = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt' | 'percentage'

export interface Operation {
  id: OperationId
  /** Short text shown on the selector button. */
  symbol: string
  /** Accessible name of the selector button. */
  label: string
  /** Labels of the operand fields; their count is the operation's arity. */
  inputs: readonly [string] | readonly [string, string]
  /** Renders the expression shown above the result. */
  expression: (a: string, b: string) => string
}

const NUMBERS = ['First number', 'Second number'] as const

/** Single source for the selector, the operand labels, the accessible names and the API call. */
export const OPERATIONS: readonly Operation[] = [
  { id: 'add', symbol: '+', label: 'Add', inputs: NUMBERS, expression: (a, b) => `${a} + ${b}` },
  { id: 'subtract', symbol: '−', label: 'Subtract', inputs: NUMBERS, expression: (a, b) => `${a} − ${b}` },
  { id: 'multiply', symbol: '×', label: 'Multiply', inputs: NUMBERS, expression: (a, b) => `${a} × ${b}` },
  { id: 'divide', symbol: '÷', label: 'Divide', inputs: NUMBERS, expression: (a, b) => `${a} ÷ ${b}` },
  {
    id: 'power',
    symbol: 'xʸ',
    label: 'Power',
    inputs: ['Base', 'Exponent'],
    expression: (a, b) => `${a} ^ ${b}`,
  },
  { id: 'sqrt', symbol: '√x', label: 'Square root', inputs: ['Number'], expression: (a) => `√${a}` },
  {
    id: 'percentage',
    symbol: 'x% of y',
    label: 'Percentage: x% of y',
    inputs: ['Percent', 'Value'],
    expression: (a, b) => `${a}% of ${b}`,
  },
]

export function getOperation(id: OperationId): Operation {
  const op = OPERATIONS.find((o) => o.id === id)
  if (!op) throw new Error(`unknown operation: ${id}`)
  return op
}
