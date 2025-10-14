import { ClientConfig } from '..'
import { User } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = User[]

const userSearch = async (
  key: string,
  config: ClientConfig,
  email: string,
  orgId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/users/search`

  const query = new URLSearchParams({
    q: email,
  })
  if (orgId) query.set('orgId', orgId)

  return await apiRequest<ResponseData>(url, 'GET', null, { authorization: key, query })
}

export default userSearch
