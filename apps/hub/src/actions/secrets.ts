'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddSecret = async (spaceId: string, name: string, value: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the space.',
    }

  const { errors } = await gruentClient(accessToken).secrets.create(spaceId, name, value)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath(`/organizations/[orgId]/spaces/[spaceId]/secrets`, 'page')
  return { error: null }
}

export const handleDeleteSecret = async (spaceId: string, secretName: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the space.',
    }

  const { errors } = await gruentClient(accessToken).secrets.delete(spaceId, secretName)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath(`/organizations/[orgId]/spaces/[spaceId]/secrets`, 'page')
  return { error: null }
}

export const handleUpdateSecret = async (spaceId: string, name: string, value: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the space.',
    }

  const { errors } = await gruentClient(accessToken).secrets.update(spaceId, name, value)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath(`/organizations/[orgId]/spaces/[spaceId]/secrets`, 'page')
  return { error: null }
}
