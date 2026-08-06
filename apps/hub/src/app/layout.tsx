import ToastNotificationContainer from '@gruent/libs/frontend/components/templates/ToastNotificationContainer'
import './styles.css'
import SettingsLoader from '@gruent/hub/components/templates/SettingsLoader'
import Providers from '@gruent/hub/contexts/Providers'
import { cookies } from 'next/headers'
export const dynamic = 'force-dynamic'

export const metadata = {
  title: 'Hub | GRUENT',
}

const RootLayout = async ({ children }: { children: React.ReactNode }) => {
  const cookieStore = await cookies()
  const accessToken = cookieStore.get('access_token')?.value

  return (
    <html>
      <head />
      <body>
        <div className="flex flex-col min-h-screen">
          <Providers accessToken={accessToken ?? null}>
            <ToastNotificationContainer />
            <SettingsLoader />
            {children}
          </Providers>
        </div>
      </body>
    </html>
  )
}

export default RootLayout
