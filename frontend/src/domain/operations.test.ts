import { describe, expect, it } from 'vitest'

import { getOperation, OPERATIONS, type OperationId } from './operations'

describe('operations', () => {
  it.each<[OperationId, string]>([
    ['add', '10 + 4'],
    ['subtract', '10 − 4'],
    ['multiply', '10 × 4'],
    ['divide', '10 ÷ 4'],
    ['power', '10 ^ 4'],
    ['sqrt', '√10'],
    ['percentage', '10% of 4'],
  ])('renders the %s expression', (id, expected) => {
    expect(getOperation(id).expression('10', '4')).toBe(expected)
  })

  it('has unique ids', () => {
    expect(new Set(OPERATIONS.map((op) => op.id)).size).toBe(OPERATIONS.length)
  })

  it('only square root is unary', () => {
    expect(OPERATIONS.filter((op) => op.inputs.length === 1).map((op) => op.id)).toEqual(['sqrt'])
  })

  it('throws for an unknown id', () => {
    expect(() => getOperation('modulo' as OperationId)).toThrow('unknown operation')
  })
})
