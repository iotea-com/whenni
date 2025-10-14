import { render } from '@testing-library/react'

import ToastNotification from '.'

describe('ToastNotification Test', () => {
  it('should render successfully', () => {
    const { baseElement } = render(
      <ToastNotification title={'Test'} body={'test'} keyName={'test'} />,
    )
    expect(baseElement).toBeTruthy()
  })
})
