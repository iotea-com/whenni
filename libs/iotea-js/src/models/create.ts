import { ClientConfig } from '..'
import { Model } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'
import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'

export type CreateModelInput = Pick<Model, 'name'> & {
  attributes: ModelAttributes
}

type ResponseData = Model

const createModel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  modelInput: CreateModelInput,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/models`

  const requestBody = modelInput

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default createModel
