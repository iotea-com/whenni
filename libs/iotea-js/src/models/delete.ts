import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteModel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  modelId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/models/${modelId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'DELETE', null, {
    authorization: key,
    query: queryParams,
  })
}

export default deleteModel
