import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { decodeJwt } from 'jose'
import { refresh, signout } from '@gruent/hub/actions/auth'
import { useCallback, useEffect, useState } from 'react'
import { useAuthProvider } from '@gruent/hub/contexts/Auth'

const useAuth = () => {
  const { accessToken } = useAuthProvider()

  const [userId, setUserId] = useState<string | null>(null)

  const getSessionFromToken = useCallback((token: string) => {
    const claims = decodeJwt(token)
    const userId = claims.sub as string
    const expiresAt = new Date((claims.exp as number) * 1000)

    const isExpired = expiresAt < new Date()
    const isExpiringSoon = expiresAt < new Date(Date.now() + 1000 * 60 * 3)
    const shouldRefresh = isExpired || isExpiringSoon

    return { userId, expiresAt, shouldRefresh }
  }, [])

  const handleRefreshToken = useCallback(async (expiresAt: Date) => {
    const isExpired = expiresAt < new Date()
    const isExpiringSoon = expiresAt < new Date(Date.now() + 1000 * 60 * 3)
    const shouldRefresh = isExpired || isExpiringSoon

    if (!shouldRefresh) return

    const { error } = await refresh()
    if (error) {
      if (!error.includes('Refresh token expired')) {
        addToast({
          title: 'Error refreshing session',
          body: error,
          level: 'error',
        })
      }

      await signout()
    }
  }, [])

  // Check if the token should be refreshed and set the userId
  useEffect(() => {
    if (!accessToken) return

    const session = getSessionFromToken(accessToken)
    setUserId(session.userId)

    // Check if the token should be refreshed
    if (session.shouldRefresh) handleRefreshToken(session.expiresAt)

    // Create a timer to check if the token should be refreshed every 3 minutes
    const refreshFunction = () => handleRefreshToken(session.expiresAt)
    const refreshInterval = 1000 * 60 * 3 // 3 minutes
    const refreshTimer = setInterval(refreshFunction, refreshInterval)

    // Add event listeners to handle focus, visibility, and network changes
    window.addEventListener('focus', refreshFunction)
    document.addEventListener('visibilitychange', refreshFunction)
    window.addEventListener('online', refreshFunction)

    return () => {
      window.removeEventListener('focus', refreshFunction)
      document.removeEventListener('visibilitychange', refreshFunction)
      window.removeEventListener('online', refreshFunction)
      clearInterval(refreshTimer)
    }
  }, [accessToken, handleRefreshToken, getSessionFromToken])

  return { accessToken, userId }
}

export default useAuth
