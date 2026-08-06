import { ClientConfig } from '..'
import { Thing } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Thing

const updateThing = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  thingId: string,
  thingInput: Pick<Thing, 'name' | 'attributes'>,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/things/${thingId}`

  const requestBody = { name: thingInput.name, attributes: thingInput.attributes }

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PUT', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default updateThing
