import { Thing } from "@prisma/client"

export type KafkaMessageQueueSourceNodeConfig = {
  "kafkaConsumer::thing"?: string | Thing
  topic?: string
}

