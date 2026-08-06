import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const healthcheck = async (config: ClientConfig): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = 'healthcheck'

  return await apiRequest<ResponseData>(url, 'GET', null)
}

export default healthcheck
