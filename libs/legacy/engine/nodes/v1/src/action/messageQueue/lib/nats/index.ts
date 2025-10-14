import { Thing } from "@prisma/client"

export type NatsMessageQueueActionNodeConfig = {
  "natsClient::thing"?: string | Thing
  topic?: string
}

