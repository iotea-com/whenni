import { ClientConfig } from '..'
import { ApiKey } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = ApiKey

const addApiKey = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  permissionSetId: string,
  name?: string,
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/api-keys`

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(
    url,
    'POST',
    { permissionSetId, name },
    { authorization: key, query: queryParams },
  )
}

export default addApiKey
