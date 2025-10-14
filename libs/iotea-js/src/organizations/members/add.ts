import { ClientConfig } from '../..'
import apiRequest from '../../request'
import { ApiResponse } from '../../response'

type ResponseData = null

const addMember = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  userId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/members`

  const requestBody = { userId }

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default addMember
