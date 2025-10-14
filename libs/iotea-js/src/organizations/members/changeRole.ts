import { ClientConfig } from '../..'
import apiRequest from '../../request'
import { ApiResponse } from '../../response'

type ResponseData = null

const changeMemberRole = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  userId: string,
  role: 'ADMIN' | 'MEMBER',
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/members/changeRole`

  const requestBody = { userId, role }

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'PATCH', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default changeMemberRole
