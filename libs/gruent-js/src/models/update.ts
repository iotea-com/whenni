import { ClientConfig } from '..'
import { Model } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Model

const updateModel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  model: Model,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/models/${model.id}`

  const requestBody = { model }

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PUT', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default updateModel
