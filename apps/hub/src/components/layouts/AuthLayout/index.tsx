import getAccessToken, { getSession } from '@gruent/hub/util/getAccessToken'
import { redirect } from 'next/navigation'
// import CallbackUrlRedirect from './CallbackUrlRedirect'

const AuthLayout = async ({ children }) => {
  // const session = await auth()
  const accessToken = await getAccessToken()
  const { userId } = await getSession(accessToken)

  // Redirect to dashboard if user is already logged in
  if (accessToken && userId) redirect('/dashboard')

  return (
    <div className="h-screen w-screen">
      {/* <CallbackUrlRedirect session={session} /> */}
      <div
        className="flex items-center justify-center h-full"
        style={{ background: 'linear-gradient(-15deg, #063712, #34dd4c)' }}
      >
        {children}
      </div>
    </div>
  )
}

export default AuthLayout
