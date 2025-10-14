import healthcheck from '.'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

describe('healthcheck', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return success when healthcheck is OK', async () => {
    const mockResponseCode = 200
    const mockResponseBody = { data: null, errors: null }
    createApiCallMock<null>(mockResponseCode, mockResponseBody)

    // When
    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await healthcheck(config)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/healthcheck')

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      method: 'GET',
    })
    expect(result).toEqual({ data: null, errors: null })
  })

  it('should return errors when healthcheck fails', async () => {
    const mockResponseCode = 500
    const mockResponseBody = 'Internal Server Error'
    createApiCallMock<null>(mockResponseCode, mockResponseBody)

    // When
    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await healthcheck(config)

    // Then
    expect(result).toEqual({ data: null, errors: ['Internal Server Error'] })
  })
})
