import { ClientConfig } from '..'
import { PermissionSet } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = PermissionSet[]

const listPermissionSets = async (
  key: string,
  config: ClientConfig,
  scopeId: string,
  options?: ListRequestOptions,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/permissions`
  url.searchParams.set('page', options?.page?.toString() ?? '1')
  url.searchParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')

  const scopeIdType = (() => {
    if (scopeId.startsWith('o')) return 'organization'
    if (scopeId.startsWith('s')) return 'space'
    throw new Error('Invalid organization or space ID')
  })()

  const queryParams = new URLSearchParams()
  if (scopeIdType === 'organization') queryParams.set('orgId', scopeId)
  if (scopeIdType === 'space') queryParams.set('spaceId', scopeId)

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listPermissionSets
