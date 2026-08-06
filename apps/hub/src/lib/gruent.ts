import { ClientOptions, createClient } from '@gruent/libs/gruent-js/src'

const envOptions: ClientOptions = {
  url: process.env.NEXT_PUBLIC_GRUENT_URL,
}

const gruentClient = (key: string, options?: ClientOptions) =>
  createClient(key, { ...envOptions, ...options })

export default gruentClient
