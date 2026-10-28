import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiError, calculate } from './calculatorApi'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
  fetchMock.mockReset()
})

function jsonResponse(status: number, body: unknown) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

async function expectApiError(promise: Promise<unknown>, status: number, code: string) {
  const err = await promise.then(
    () => null,
    (e: unknown) => e,
  )
  expect(err).toBeInstanceOf(ApiError)
  expect(err).toMatchObject({ status, code })
  return err as ApiError
}

describe('calculate', () => {
  it('posts JSON to the operation endpoint and returns the result', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { operation: 'divide', result: 2.5 }))
    const signal = new AbortController().signal

    await expect(calculate('divide', { a: 10, b: 4 }, { signal })).resolves.toBe(2.5)

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/divide', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: '{"a":10,"b":4}',
      signal,
    })
  })

  it('sends only "a" for unary operations', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { operation: 'sqrt', result: 4 }))

    await calculate('sqrt', { a: 16 })

    expect(fetchMock.mock.calls[0][1]?.body).toBe('{"a":16}')
  })

  it('throws the API error for error responses', async () => {
    fetchMock.mockResolvedValue(
      jsonResponse(422, { error: { code: 'DIVISION_BY_ZERO', message: 'division by zero is undefined' } }),
    )

    const err = await expectApiError(calculate('divide', { a: 1, b: 0 }), 422, 'DIVISION_BY_ZERO')
    expect(err.message).toBe('division by zero is undefined')
  })

  it('handles responses that are not JSON', async () => {
    fetchMock.mockResolvedValue(new Response('<html>Bad Gateway</html>', { status: 502 }))

    await expectApiError(calculate('add', { a: 1, b: 2 }), 502, 'UNEXPECTED_RESPONSE')
  })

  it('handles a 200 without a numeric result', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { result: 'NaN' }))

    await expectApiError(calculate('add', { a: 1, b: 2 }), 200, 'UNEXPECTED_RESPONSE')
  })

  it('reports network failures', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))

    await expectApiError(calculate('add', { a: 1, b: 2 }), 0, 'NETWORK_ERROR')
  })

  it('reports timeouts', async () => {
    fetchMock.mockRejectedValue(new DOMException('signal timed out', 'TimeoutError'))

    await expectApiError(calculate('add', { a: 1, b: 2 }), 0, 'TIMEOUT')
  })
})
