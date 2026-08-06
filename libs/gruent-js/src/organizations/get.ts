import { ClientConfig } from '..'
import { Organization, Space } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Organization & { spaces: Space[] }

const getOrganizationById = async (
  key: string,
  config: ClientConfig,
  orgId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/${orgId}`

  return await apiRequest(url, 'GET', null, { authorization: key })
}

export default getOrganizationById
