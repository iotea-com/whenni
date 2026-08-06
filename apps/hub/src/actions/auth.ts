'use server'

import { cookies } from 'next/headers'
import gruentClient from '../lib/gruent'
import { redirect } from 'next/navigation'

export async function signInWithCredentials({
  email,
  password,
  redirectTo,
}: {
  email: string
  password: string
  redirectTo: string
}): Promise<{ error: string | null; redirect?: string }> {
  const { data: signinResponse, errors } = await gruentClient('').auth.signin.credentials(
    email,
    password,
    {
      redirectTo,
    },
  )

  if (errors && errors.length > 0) return { error: errors[0] }
  if (!signinResponse) return { error: 'Invalid response from server. Please try again.' }

  const cookieStore = await cookies()
  cookieStore.set('access_token', signinResponse.accessToken, {
    httpOnly: false,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 7, // 7 days
  })

  cookieStore.set('refresh_token', signinResponse.refreshToken, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 7, // 7 days
  })

  redirect(redirectTo)
}

export async function signInWithMagicLink({
  email,
  appUrl,
  redirectTo,
}: {
  email: string
  appUrl: string
  redirectTo: string
}): Promise<{ error: string | null; redirect?: string }> {
  const { errors } = await gruentClient('').auth.signin.magicLink.send(email, appUrl, { redirectTo })

  if (errors && errors.length > 0) return { error: errors[0] }

  return { error: null }
}

export async function signup({
  email,
  password,
  inviteToken,
  origin,
}: {
  email: string
  password: string
  inviteToken?: string
  origin: string
}) {
  const { errors } = await gruentClient('').auth.signup(email, password, origin, inviteToken)

  return { error: errors?.[0] ?? null }
}

export async function refresh() {
  const cookieStore = await cookies()
  const refreshToken = cookieStore.get('refresh_token')?.value

  if (!refreshToken) return { error: 'No refresh token found' }

  const { data: refreshResponse, errors } = await gruentClient('').auth.refresh(refreshToken)
  if (errors && errors.length > 0) return { error: errors[0] }

  if (!refreshResponse) return { error: 'Invalid response from server. Please try again.' }

  cookieStore.set('access_token', refreshResponse.accessToken, {
    httpOnly: false,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 7, // 7 days
  })

  cookieStore.set('refresh_token', refreshResponse.refreshToken, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 7, // 7 days
  })

  return { error: null }
}

export async function signout() {
  const cookieStore = await cookies()
  cookieStore.delete('access_token')
  cookieStore.delete('refresh_token')

  redirect('/signin')
}
