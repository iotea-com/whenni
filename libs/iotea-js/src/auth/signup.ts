import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = null

const signup = async (
  key: string,
  config: ClientConfig,
  email: string,
  password: string,
  origin: string,
  inviteToken?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/signup`

  const requestBody = {
    email,
    password,
    appUrl: origin,
    inviteToken,
  }

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
  })
}

export default signup
