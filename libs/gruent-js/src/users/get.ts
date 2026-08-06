import { ClientConfig } from '..'
import { Organization, OrganizationMember, Space, User } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = User & {
  organizations: (OrganizationMember & {
    organization: Organization & { spaces: Space[] }
  })[]
}

const getUserById = async (
  key: string,
  config: ClientConfig,
  userId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/users/${userId}`

  return await apiRequest<ResponseData>(url, 'GET', null, { authorization: key })
}

export default getUserById
