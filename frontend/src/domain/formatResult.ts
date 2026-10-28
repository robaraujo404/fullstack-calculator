/**
 * Formats an API result for display. The value returned by the API is not
 * changed; only its presentation is.
 *
 * Integers up to 2^53 are shown exactly: rounding them to 15 digits would turn
 * 9007199254740991 into 9007199254740990. Larger integers use scientific
 * notation, because String() pads their shortest round-trip digits with zeros
 * (2^60 would read 1152921504606847000 instead of ...976). Other values are
 * rounded to 15 significant digits, which hides float64 noise such as
 * 0.1 + 0.2 = 0.30000000000000004.
 */
export function formatResult(n: number): string {
  if (!Number.isInteger(n)) return String(Number(n.toPrecision(15)))
  return Number.isSafeInteger(n) ? String(n) : n.toExponential()
}
