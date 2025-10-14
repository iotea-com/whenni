import { ClientConfig } from '..'
import { PermissionSet } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = PermissionSet

const updatePermissionSet = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  permissionSetId: string,
  name: string,
  permissions: string[],
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/permissions/${permissionSetId}`

  const requestBody = { name, permissions }

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PUT', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default updatePermissionSet
