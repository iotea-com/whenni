'use server'

import ioteaClient from '@iotea/hub/lib/iotea'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@iotea/hub/util/getAccessToken'

type MemberInvitation = {
  orgId: string
  email: string
  invitedBy: string
}

export const handleAddMember = async (orgId: string, userId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the member.',
    }

  const { errors } = await ioteaClient(accessToken).organizations.members.add(orgId, userId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('organizations/[orgId]/team', 'page')
  revalidatePath('organizations/[orgId]/spaces/[spaceId]/team', 'page')
  return { error: null }
}

export const handleRemoveMember = async (orgId: string, userId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to remove the member.',
    }

  const { errors } = await ioteaClient(accessToken).organizations.members.remove(orgId, userId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('organizations/[orgId]/team', 'page')
  revalidatePath('organizations/[orgId]/spaces/[spaceId]/team', 'page')
  return { error: null }
}

export const handleChangeMemberRole = async (
  orgId: string,
  userId: string,
  role: 'ADMIN' | 'MEMBER',
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to change the member role.',
    }

  const { errors } = await ioteaClient(accessToken).organizations.members.changeRole(
    orgId,
    userId,
    role,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('organizations/[orgId]/team', 'page')
  revalidatePath('organizations/[orgId]/spaces/[spaceId]/team', 'page')
  return { error: null }
}

export const handleInviteMember = async (orgId: string, email: string, origin: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to invite a member.',
    }

  const { errors } = await ioteaClient(accessToken).organizations.members.invite(
    orgId,
    email,
    origin,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  return { error: null }
}

export const verifyInvitationToken = async (secret: string, token: string) => {
  try {
    // Split the token into its components
    const [encodedHeader, encodedPayload, receivedSignature] = token.split('.')

    // Verify signature
    const enc = new TextEncoder()
    const key = await crypto.subtle.importKey(
      'raw',
      enc.encode(secret),
      { name: 'HMAC', hash: 'SHA-256' },
      false,
      ['verify'],
    )

    const isValid = await crypto.subtle.verify(
      'HMAC',
      key,
      Buffer.from(receivedSignature, 'base64url'),
      enc.encode(`${encodedHeader}.${encodedPayload}`),
    )

    if (!isValid)
      return {
        error: 'Invalid signature',
      }

    // Decode and verify payload
    const payload = JSON.parse(
      Buffer.from(encodedPayload, 'base64url').toString(),
    ) as MemberInvitation & {
      exp: number
      iss: string
    }

    // Check expiration
    if (payload.exp < Math.floor(Date.now() / 1000))
      return {
        error: 'Token has expired',
      }

    return {
      valid: true,
      data: payload,
    }
  } catch (err) {
    return {
      error: 'Invalid or expired invitation token',
      originalError: err,
    }
  }
}
