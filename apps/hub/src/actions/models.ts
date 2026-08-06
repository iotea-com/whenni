'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { CreateModelInput } from '@gruent/libs/gruent-js/src/models/create'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddModel = async (spaceId: string, modelInput: CreateModelInput) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the model.',
    }

  const { errors } = await gruentClient(accessToken).models.create(spaceId, modelInput)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/models', 'page')
  return { error: null }
}

export const handleDeleteModel = async (spaceId: string, modelId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the model.',
    }

  const { errors } = await gruentClient(accessToken).models.delete(spaceId, modelId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/models', 'page')
  return { error: null }
}
