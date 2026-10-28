import { describe, expect, it } from 'vitest'

import { formatResult } from './formatResult'

describe('formatResult', () => {
  it.each([
    [0.1 + 0.2, '0.3'],
    [2.5, '2.5'],
    [-2.5e-10, '-2.5e-10'],
    [1 / 3, '0.333333333333333'],
    [1e-7, '1e-7'],
    [1e21, '1e+21'],
    [42, '42'],
    [-7, '-7'],
    [0, '0'],
  ])('formats %d as %s', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })

  it('shows integers exactly instead of rounding them', () => {
    expect(formatResult(2 ** 53 - 1)).toBe('9007199254740991')
  })

  it.each([
    [2 ** 60, '1.152921504606847e+18'],
    [-(2 ** 60), '-1.152921504606847e+18'],
  ])('shows %d in scientific notation instead of padding it with zeros', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })

  it.each([
    [100000000000000.4, '100000000000000'],
    [123450000000000.25, '123450000000000'],
  ])('keeps the integer part of %d', (value, expected) => {
    expect(formatResult(value)).toBe(expected)
  })
})
