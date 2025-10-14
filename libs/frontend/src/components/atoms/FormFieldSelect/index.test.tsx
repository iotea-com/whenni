import { render } from '@testing-library/react'

import FormFieldSelect from '.'

describe('FormFieldSelect Test', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<FormFieldSelect name={'test'} options={[]} label={'Test'} />)
    expect(baseElement).toBeTruthy()
  })
})
