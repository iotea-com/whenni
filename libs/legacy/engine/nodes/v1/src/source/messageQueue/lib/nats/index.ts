import { Thing } from "@prisma/client"

export type NatsMessageQueueSourceNodeConfig = {
  "natsClient::thing"?: string | Thing
  topic?: string
}

