import styles from './index.module.scss'
import getAccessToken, { getSession } from '@iotea/hub/util/getAccessToken'
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
      <div className={styles.authFormWrapper}>{children}</div>
    </div>
  )
}

export default AuthLayout
