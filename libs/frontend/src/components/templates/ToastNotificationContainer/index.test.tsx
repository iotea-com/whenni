import { render } from '@testing-library/react'

import ToastNotificationContainer from '.'

describe('ToastNotificationContainer Test', () => {
  it('should render successfully', () => {
    const { baseElement } = render(<ToastNotificationContainer />)
    expect(baseElement).toBeTruthy()
  })
})
