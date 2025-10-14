import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = null

const signinMagicLink = async (
  key: string,
  config: ClientConfig,
  email: string,
  appUrl: string,
  options?: { redirectTo?: string },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/signin`

  const requestBody = {
    method: 'magicLink',
    email,
    appUrl,
    redirectTo: options?.redirectTo,
  }

  return await apiRequest<ResponseData>(url, 'POST', requestBody, { authorization: key })
}

export default signinMagicLink
