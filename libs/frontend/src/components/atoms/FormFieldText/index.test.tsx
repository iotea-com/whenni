import { render } from '@testing-library/react'

import FormFieldText from '.'

describe('FormFieldTextTest', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<FormFieldText name="test" label="Test" />)
    expect(baseElement).toBeTruthy()
  })
})
