'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { FC, PropsWithChildren } from 'react'
import AuthProvider from './Auth'

type Props = {
  accessToken: string | null
}

const Providers: FC<PropsWithChildren<Props>> = ({ children, accessToken }) => {
  const queryClient = new QueryClient()

  return (
    <AuthProvider accessToken={accessToken}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </AuthProvider>
  )
}

export default Providers
