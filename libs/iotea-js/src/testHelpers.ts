import { ApiResponse } from './response'

export function createApiCallMock<ResponseData>(
  status: number,
  responseBody: string | ApiResponse<ResponseData>,
) {
  global.fetch = jest.fn(() =>
    Promise.resolve({
      status,
      json: () => Promise.resolve(responseBody),
      text: () => Promise.resolve(responseBody),
      headers: new Headers(),
    }),
  ) as jest.Mock
}
