import { act, renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { ApiError, type CalculateFn } from '../api/calculatorApi'
import { useCalculator } from './useCalculator'

function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: unknown) => void
  const promise = new Promise<T>((res, rej) => {
    resolve = res
    reject = rej
  })
  return { promise, resolve, reject }
}

describe('useCalculator', () => {
  it('goes idle → loading → success', async () => {
    const pending = deferred<number>()
    const client = vi.fn<CalculateFn>(() => pending.promise)
    const { result } = renderHook(() => useCalculator(client))

    expect(result.current.state).toEqual({ status: 'idle' })

    let run!: Promise<void>
    act(() => {
      run = result.current.run('add', { a: 2, b: 3 })
    })
    expect(result.current.state).toEqual({ status: 'loading' })

    await act(async () => {
      pending.resolve(5)
      await run
    })
    expect(result.current.state).toEqual({ status: 'success', value: 5 })
  })

  it('passes the operands and a timeout signal to the client', async () => {
    const client = vi.fn<CalculateFn>().mockResolvedValue(4)
    const { result } = renderHook(() => useCalculator(client))

    await act(() => result.current.run('sqrt', { a: 16 }))

    expect(client).toHaveBeenCalledWith('sqrt', { a: 16 }, { signal: expect.any(AbortSignal) })
  })

  it('stores API errors', async () => {
    const client = vi.fn<CalculateFn>().mockRejectedValue(
      new ApiError(422, 'DIVISION_BY_ZERO', 'division by zero is undefined'),
    )
    const { result } = renderHook(() => useCalculator(client))

    await act(() => result.current.run('divide', { a: 1, b: 0 }))

    expect(result.current.state).toEqual({
      status: 'error',
      code: 'DIVISION_BY_ZERO',
      message: 'division by zero is undefined',
    })
  })

  it('turns unexpected errors into a generic message', async () => {
    const client = vi.fn<CalculateFn>().mockRejectedValue(new Error('boom'))
    const { result } = renderHook(() => useCalculator(client))

    await act(() => result.current.run('add', { a: 1, b: 2 }))

    expect(result.current.state).toMatchObject({ status: 'error', code: 'UNEXPECTED_ERROR' })
  })

  it('resets to idle', async () => {
    const client = vi.fn<CalculateFn>().mockResolvedValue(5)
    const { result } = renderHook(() => useCalculator(client))

    await act(() => result.current.run('add', { a: 2, b: 3 }))
    act(() => result.current.reset())

    expect(result.current.state).toEqual({ status: 'idle' })
  })
})
