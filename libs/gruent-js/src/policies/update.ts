import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import { Policy } from '.'
import apiRequest from '../request'

type ResponseData = {
  policy: Policy
  revoke: boolean
}

const updatePolicy = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  certificateId: string,
  policy: Policy,
  revoke: boolean,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/policies/${certificateId}`

  const requestBody = { policy, revoke }

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PUT', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default updatePolicy
