'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddSpace = async (orgId: string, name: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the space.',
    }

  const { errors } = await gruentClient(accessToken).spaces.create(orgId, name)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/dashboard', 'page')
  return { error: null }
}

export const handleDeleteSpace = async (orgId: string, spaceId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the space.',
    }

  const { errors } = await gruentClient(accessToken).spaces.delete(orgId, spaceId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/dashboard', 'page')
  return { error: null }
}
