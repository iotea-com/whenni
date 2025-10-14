import { ClientConfig } from '..'
import { PermissionSet } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = PermissionSet

const addPermissionSet = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  name: string,
  permissions: string[],
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/permissions`

  const requestBody: Record<string, any> = { name, permissions }

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default addPermissionSet
