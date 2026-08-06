import { ClientConfig } from '..'
import { Organization } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Organization

const createOrganization = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  name: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations`

  const requestBody = { name }

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default createOrganization
