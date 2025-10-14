import { jwtVerify } from 'jose'
import { cookies } from 'next/headers'

const getAccessToken = async () => {
  const cookieStore = await cookies()
  const accessToken = cookieStore.get('access_token')

  return accessToken?.value
}

const getSession = async (accessToken: string | undefined) => {
  if (!accessToken) return { userId: null, expiresAt: null }

  try {
    const { payload } = await jwtVerify(
      accessToken,
      new TextEncoder().encode(process.env.JWT_SECRET),
      {
        algorithms: ['HS256'],
      },
    )
    const userId = payload.sub as string
    const expiresAt = new Date((payload.exp as number) * 1000)

    return { userId, expiresAt }
  } catch (_error) {
    return { userId: null, expiresAt: null }
  }
}

export default getAccessToken
export { getSession }
