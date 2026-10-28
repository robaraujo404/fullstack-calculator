export type ParseResult = { ok: true; value: number } | { ok: false; error: string }

// Optional sign, digits with an optional decimal point, optional exponent.
// Validating the whole string first avoids Number('') === 0,
// Number('0x10') === 16 and parseFloat('1abc') === 1.
const NUMBER = /^[+-]?(\d+\.?\d*|\.\d+)([eE][+-]?\d+)?$/
const DECIMAL_COMMA = /^[+-]?\d*,\d+$/

export function parseOperand(raw: string): ParseResult {
  const text = raw.trim()

  if (text === '') {
    return { ok: false, error: 'Enter a number.' }
  }
  if (DECIMAL_COMMA.test(text)) {
    return { ok: false, error: 'Use a dot as the decimal separator, e.g. 1.5.' }
  }
  if (!NUMBER.test(text)) {
    return { ok: false, error: 'Enter a valid number, e.g. 42, -3.5 or 1e6.' }
  }

  const value = Number(text)
  if (!Number.isFinite(value)) {
    return { ok: false, error: 'This number is too large.' }
  }
  return { ok: true, value }
}
