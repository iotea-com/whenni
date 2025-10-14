import ioteaClient from '@iotea/hub/lib/iotea'
import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { NextRequest, NextResponse } from 'next/server'

const handler = async (req: NextRequest) => {
  const { searchParams } = new URL(req.url)
  const { t: verificationToken } = Object.fromEntries(searchParams.entries())

  // Verify magic link
  const { data: verificationResponse, errors: tokenErrors } =
    await ioteaClient('').auth.signin.magicLink.verify(verificationToken)

  if (tokenErrors) redirect(`/magic-link/verify/error?error=${encodeURIComponent(tokenErrors[0])}`)
  if (!verificationResponse)
    redirect(`/magic-link/verify/error?error=${encodeURIComponent('Invalid magic link')}`)

  // Set auth cookie
  const cookieStore = await cookies()
  cookieStore.set('access_token', verificationResponse.accessToken, {
    httpOnly: false,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 30, // 30 days
  })

  cookieStore.set('refresh_token', verificationResponse.refreshToken, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge: 60 * 60 * 24 * 30, // 30 days
  })

  return NextResponse.redirect(new URL('/dashboard', req.url))
}

export { handler as GET }
