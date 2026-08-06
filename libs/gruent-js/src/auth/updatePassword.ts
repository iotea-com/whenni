import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = null

const changePassword = async (
  key: string,
  config: ClientConfig,
  password: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/auth/password/update`

  const requestBody = {
    password,
  }

  return await apiRequest<ResponseData>(url, 'PATCH', requestBody, { authorization: key })
}

export default changePassword
