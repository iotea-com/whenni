import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = {
  accessToken: string
  refreshToken: string
}

const signinCredentials = async (
  key: string,
  config: ClientConfig,
  email: string,
  password: string,
  options?: { redirectTo?: string },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/signin`

  const requestBody = {
    method: 'credentials',
    email,
    password,
    redirectTo: options?.redirectTo,
  }

  return await apiRequest<ResponseData>(url, 'POST', requestBody, { authorization: key })
}

export default signinCredentials
