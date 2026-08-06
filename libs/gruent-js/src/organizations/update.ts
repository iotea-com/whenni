import { ClientConfig } from '..'
import { Organization, Space } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Space

const updateOrganization = async (
  key: string,
  config: ClientConfig,
  organization: Organization,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/${organization.id}`

  return await apiRequest<ResponseData>(url, 'PUT', { organization }, { authorization: key })
}

export default updateOrganization
