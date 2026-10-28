import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { ApiError, type CalculateFn } from '../api/calculatorApi'
import { Calculator } from './Calculator'

function setup(client: CalculateFn = vi.fn<CalculateFn>().mockResolvedValue(0)) {
  const user = userEvent.setup()
  render(<Calculator client={client} />)
  return {
    user,
    client,
    first: () => screen.getByRole('textbox', { name: /first number|^number/i }),
    second: () => screen.queryByRole('textbox', { name: 'Second number' }),
    calculate: () => screen.getByRole('button', { name: 'Calculate' }),
    result: () => screen.getByRole('region', { name: 'Result' }),
  }
}

describe('Calculator', () => {
  it('sends the operands to the selected operation and shows the result', async () => {
    const client = vi.fn<CalculateFn>().mockResolvedValue(2.5)
    const { user, first, second, calculate, result } = setup(client)

    await user.type(first(), '10')
    await user.click(screen.getByRole('radio', { name: 'Divide' }))
    await user.type(second()!, '4')
    await user.click(calculate())

    expect(client).toHaveBeenCalledWith('divide', { a: 10, b: 4 }, { signal: expect.any(AbortSignal) })
    expect(await screen.findByText('2.5')).toBeInTheDocument()
    expect(result()).toHaveTextContent('10 ÷ 4')
  })

  it('formats float noise away', async () => {
    const { user, first, second, calculate } = setup(vi.fn<CalculateFn>().mockResolvedValue(0.1 + 0.2))

    await user.type(first(), '0.1')
    await user.type(second()!, '0.2')
    await user.click(calculate())

    expect(await screen.findByText('0.3')).toBeInTheDocument()
  })

  it('blocks invalid input before calling the API', async () => {
    const { user, client, first, second, calculate } = setup()

    await user.type(first(), 'abc')
    await user.type(second()!, '1,5')
    await user.click(calculate())

    expect(client).not.toHaveBeenCalled()
    expect(first()).toHaveAttribute('aria-invalid', 'true')
    expect(first()).toHaveAccessibleDescription(/valid number/)
    expect(second()).toHaveAccessibleDescription(/dot as the decimal separator/)
  })

  it('shows API errors and moves focus to them', async () => {
    const client = vi
      .fn<CalculateFn>()
      .mockRejectedValue(new ApiError(422, 'DIVISION_BY_ZERO', 'division by zero is undefined'))
    const { user, first, second, calculate, result } = setup(client)

    await user.type(first(), '1')
    await user.click(screen.getByRole('radio', { name: 'Divide' }))
    await user.type(second()!, '0')
    await user.click(calculate())

    expect(await screen.findByRole('alert')).toHaveTextContent('Division by zero is undefined.')
    expect(result()).toHaveFocus()
  })

  it('falls back to the API message for other errors', async () => {
    const client = vi.fn<CalculateFn>().mockRejectedValue(new ApiError(0, 'NETWORK_ERROR', 'Could not reach the server.'))
    const { user, first, second, calculate } = setup(client)

    await user.type(first(), '1')
    await user.type(second()!, '2')
    await user.click(calculate())

    expect(await screen.findByRole('alert')).toHaveTextContent('Could not reach the server.')
  })

  it('hides the second input for square root and sends only "a"', async () => {
    const client = vi.fn<CalculateFn>().mockResolvedValue(4)
    const { user, first, second, calculate } = setup(client)

    await user.click(screen.getByRole('radio', { name: 'Square root' }))
    expect(second()).not.toBeInTheDocument()

    await user.type(first(), '16')
    await user.click(calculate())

    expect(client).toHaveBeenCalledWith('sqrt', { a: 16 }, expect.anything())
  })

  it('submits with Enter', async () => {
    const client = vi.fn<CalculateFn>().mockResolvedValue(5)
    const { user, first, second } = setup(client)

    await user.type(first(), '2')
    await user.type(second()!, '3{Enter}')

    expect(client).toHaveBeenCalledWith('add', { a: 2, b: 3 }, expect.anything())
  })

  it('changes the operation with the arrow keys', async () => {
    const { user } = setup()

    await user.click(screen.getByRole('radio', { name: 'Add' }))
    await user.keyboard('{ArrowRight}')

    expect(screen.getByRole('radio', { name: 'Subtract' })).toBeChecked()
  })

  it('disables the form while the request runs and focuses the result after', async () => {
    let resolve!: (value: number) => void
    const client = vi.fn<CalculateFn>(() => new Promise((res) => (resolve = res)))
    const { user, first, second, calculate, result } = setup(client)

    await user.type(first(), '2')
    await user.type(second()!, '3')
    await user.click(calculate())

    expect(first()).toBeDisabled()
    expect(calculate()).toBeDisabled()
    expect(screen.getByRole('radio', { name: 'Add' })).toBeDisabled()

    resolve(5)

    expect(await screen.findByText('5')).toBeInTheDocument()
    expect(first()).toBeEnabled()
    expect(result()).toHaveFocus()
  })

  it('clears the previous result when an input changes', async () => {
    const { user, first, second, calculate } = setup(vi.fn<CalculateFn>().mockResolvedValue(5))

    await user.type(first(), '2')
    await user.type(second()!, '3')
    await user.click(calculate())
    expect(await screen.findByText('5')).toBeInTheDocument()

    await user.type(first(), '1')

    expect(screen.queryByText('5')).not.toBeInTheDocument()
    expect(screen.queryByText('21 + 3')).not.toBeInTheDocument()
  })
})
