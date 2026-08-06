import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = {
  accessToken: string
  refreshToken: string
}

const verifyMagicLink = async (
  key: string,
  config: ClientConfig,
  token: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/magic-link/verify`

  const requestBody = { token }

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
  })
}

export default verifyMagicLink
