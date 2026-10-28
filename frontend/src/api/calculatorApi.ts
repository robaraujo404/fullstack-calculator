import type { OperationId } from '../domain/operations'

export class ApiError extends Error {
  /** HTTP status, or 0 when no response was received. */
  readonly status: number
  readonly code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

export interface Operands {
  a: number
  b?: number
}

export type CalculateFn = (
  op: OperationId,
  operands: Operands,
  options?: { signal?: AbortSignal },
) => Promise<number>

/**
 * Calls POST /api/v1/{op}. The path is relative: the Vite dev server proxies
 * it to the backend, and in Docker the backend serves this app itself.
 * Every failure is thrown as an ApiError, so callers never deal with HTTP.
 */
export const calculate: CalculateFn = async (op, operands, options = {}) => {
  let response: Response
  try {
    response = await fetch(`/api/v1/${op}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(operands),
      signal: options.signal,
    })
  } catch (err) {
    if (isTimeout(err)) {
      throw new ApiError(0, 'TIMEOUT', 'The server took too long to respond.')
    }
    throw new ApiError(0, 'NETWORK_ERROR', 'Could not reach the server. Is the backend running?')
  }

  const body: unknown = await response.json().catch(() => null)
  if (response.ok && hasResult(body)) {
    return body.result
  }
  if (hasError(body)) {
    throw new ApiError(response.status, body.error.code, body.error.message)
  }
  throw new ApiError(
    response.status,
    'UNEXPECTED_RESPONSE',
    `Unexpected response from the server (HTTP ${response.status}).`,
  )
}

// Checks the name rather than `instanceof DOMException`, which fails when the
// exception comes from another realm (as it does under jsdom).
function isTimeout(err: unknown): boolean {
  return typeof err === 'object' && err !== null && (err as { name?: unknown }).name === 'TimeoutError'
}

function hasResult(body: unknown): body is { result: number } {
  return typeof body === 'object' && body !== null && typeof (body as { result?: unknown }).result === 'number'
}

function hasError(body: unknown): body is { error: { code: string; message: string } } {
  if (typeof body !== 'object' || body === null) return false
  const error = (body as { error?: { code?: unknown; message?: unknown } }).error
  return typeof error?.code === 'string' && typeof error.message === 'string'
}
