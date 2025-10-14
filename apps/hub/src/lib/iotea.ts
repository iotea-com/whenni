import { ClientOptions, createClient } from '@iotea/libs/iotea-js/src'

const envOptions: ClientOptions = {
  url: process.env.NEXT_PUBLIC_IOTEA_URL,
}

const ioteaClient = (key: string, options?: ClientOptions) =>
  createClient(key, { ...envOptions, ...options })

export default ioteaClient
