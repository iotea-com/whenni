'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddOrganization = async (userId: string, name: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the space.',
    }

  const { errors } = await gruentClient(accessToken).organizations.create(userId, name)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/dashboard', 'page')
  return { error: null }
}

export const handleDeleteOrganization = async (orgId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the space.',
    }

  const { errors } = await gruentClient(accessToken).organizations.delete(orgId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/dashboard', 'page')
  return { error: null }
}
