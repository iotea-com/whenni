import { Thing } from "@prisma/client"

export type KafkaMessageQueueActionNodeConfig = {
  "kafkaProducer::thing"?: string | Thing
  topic?: string
}

