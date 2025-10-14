import AuthLayout from '@iotea/hub/components/layouts/AuthLayout'

export const metadata = {
  title: 'Verify email | IOTEA',
}

const Layout = async ({ children }) => {
  return <AuthLayout>{children}</AuthLayout>
}

export default Layout
