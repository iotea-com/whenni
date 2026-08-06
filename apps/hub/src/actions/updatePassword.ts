'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import getAccessToken from '@gruent/hub/util/getAccessToken'

const updatePassword = async ({ newPassword }: { newPassword: string }) => {
  if (!newPassword) return { error: 'No new password was provided.' }

  const accessToken = await getAccessToken()
  if (!accessToken) return { error: 'Not signed in.' }

  const { errors } = await gruentClient(accessToken).auth.updatePassword(newPassword)
  if (errors) return { error: errors[0] }

  return { error: null }
}

export default updatePassword
