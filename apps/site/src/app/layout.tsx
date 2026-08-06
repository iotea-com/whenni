import Navbar from '@gruent/libs/frontend/components/organisms/Navbar'
import Footer from '@gruent/libs/frontend/components/organisms/Footer'

import './styles.css'
import Providers from './Providers'
import BetaSignupModal from '@gruent/libs/frontend/components/organisms/Modal/BetaSignupModal'
import MobileNavModal from '@gruent/libs/frontend/components/organisms/Modal/MobileNavModal'

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <html>
      <head />
      <body>
        <Providers>
          <div className="flex flex-col min-h-screen">
            <BetaSignupModal />
            <MobileNavModal />
            <Navbar />
            <div className="grow">{children}</div>
            <Footer />
          </div>
        </Providers>
      </body>
    </html>
  )
}

export default RootLayout
