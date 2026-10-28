import { describe, expect, it } from 'vitest'

import { parseOperand } from './parseOperand'

describe('parseOperand', () => {
  it.each([
    ['42', 42],
    ['-3.5', -3.5],
    ['+7', 7],
    ['.5', 0.5],
    ['-0.5', -0.5],
    ['5.', 5],
    ['1e3', 1000],
    ['2.5E-3', 0.0025],
    ['  12  ', 12],
    ['0', 0],
    ['0e-400', 0],
  ])('accepts %j as %d', (raw, value) => {
    expect(parseOperand(raw)).toEqual({ ok: true, value })
  })

  it.each([
    ['', 'Enter a number'],
    ['   ', 'Enter a number'],
    ['1,5', 'dot as the decimal separator'],
    ['abc', 'valid number'],
    ['1abc', 'valid number'],
    ['0x10', 'valid number'],
    ['Infinity', 'valid number'],
    ['NaN', 'valid number'],
    ['1e', 'valid number'],
    ['1.2.3', 'valid number'],
    ['--1', 'valid number'],
    ['1e400', 'too large'],
    ['1e-400', 'too small'],
  ])('rejects %j', (raw, message) => {
    const result = parseOperand(raw)
    expect(result.ok).toBe(false)
    expect(result.ok ? '' : result.error).toContain(message)
  })
})
