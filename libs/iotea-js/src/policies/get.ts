import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import { Policy } from '.'
import apiRequest from '../request'

type ResponseData = {
  policy: Policy
  revoke: boolean
}

const getPolicies = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  certificateId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/policies/${certificateId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default getPolicies
