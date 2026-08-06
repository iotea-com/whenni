import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = {
  accessToken: string
  refreshToken: string
}

const refresh = async (
  key: string,
  config: ClientConfig,
  refreshToken: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/refresh`

  const requestBody = {
    refreshToken,
  }

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
  })
}

export default refresh
