import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteOrganizationById = async (
  key: string,
  config: ClientConfig,
  orgId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/${orgId}`

  return await apiRequest<ResponseData>(url, 'DELETE', null, { authorization: key })
}

export default deleteOrganizationById
