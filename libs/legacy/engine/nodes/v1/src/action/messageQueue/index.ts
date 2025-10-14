import { ChannelNode } from "../../../node"
import { KafkaMessageQueueActionNodeConfig } from './lib/kafka'
import { NatsMessageQueueActionNodeConfig } from './lib/nats'

export type MessageQueueActionNodeConfig = KafkaMessageQueueActionNodeConfig | NatsMessageQueueActionNodeConfig

export const defaultMessageQueueActionNode: ChannelNode<MessageQueueActionNodeConfig> = {
  id: "message-queue-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Message Queue",
    label: "messageQueue",
    type: "action",
    description: "Publish messages to a topic (queue, subject, or channel) in a message queue",
    io: {
      inputs: [
        {
          id: "messageQueueInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [],
    },
    config: {},
    dependencies: {
      things: [],
      models: [],
    },
  },
}
