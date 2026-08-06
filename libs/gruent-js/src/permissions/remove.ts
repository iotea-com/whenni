import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deletePermissionSet = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  permissionSetId: string,
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/permissions`

  const requestBody = { permissionSetId }

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'DELETE', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default deletePermissionSet
