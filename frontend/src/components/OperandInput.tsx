interface OperandInputProps {
  id: string
  label: string
  value: string
  error?: string
  onChange: (value: string) => void
}

/** A text input for one operand, with its label and inline error. */
export function OperandInput({ id, label, value, error, onChange }: OperandInputProps) {
  const errorId = `${id}-error`
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        type="text"
        // Numeric keypad on mobile without type="number", which accepts a bare
        // "e", varies across browsers and hides the raw text from validation.
        inputMode="decimal"
        autoComplete="off"
        spellCheck={false}
        placeholder="e.g. 3.14"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
      />
      {error && (
        <p id={errorId} className="field-error">
          {error}
        </p>
      )}
    </div>
  )
}
