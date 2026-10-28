/**
 * Formats an API result for display. The value returned by the API is not
 * changed; only its presentation is.
 *
 * Integers are shown exactly: rounding them to 15 digits would turn
 * 9007199254740991 into 9007199254740990. Other values are rounded to 15
 * significant digits, which hides float64 noise such as
 * 0.1 + 0.2 = 0.30000000000000004.
 */
export function formatResult(n: number): string {
  return Number.isInteger(n) ? String(n) : String(Number(n.toPrecision(15)))
}
