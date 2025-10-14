import { ChannelNode } from "../../../node"
import { KafkaMessageQueueSourceNodeConfig } from './lib/kafka'
import { NatsMessageQueueSourceNodeConfig } from './lib/nats'

export type MessageQueueSourceNodeConfig = KafkaMessageQueueSourceNodeConfig | NatsMessageQueueSourceNodeConfig

export const defaultMessageQueueSourceNode: ChannelNode<MessageQueueSourceNodeConfig> = {
  id: "message-queue-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Message Queue",
    label: "messageQueue",
    type: "source",
    description: "Subscribe to messages from a topic (queue, subject, or channel) in a message queue",
    io: {
      inputs: [
        {
          id: "passThrough",
          allowedEdge: ["channel"],
          dataType: ["bytes"],
        }
      ],
      outputs: [
        {
          id: 'trigger',
          allowedEdge: ["processing", "action", "conditional", "custom"],
          dataType: ["signal"],
        }
      ],
    },
    config: {},
    dependencies: {
      things: [],
      models: [],
    },
  },
}
