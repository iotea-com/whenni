import { ChannelNode } from "../../../node"

export type TimerSourceConfig = {
  cronExpression: string
}

export const defaultTimerSourceNode: ChannelNode<TimerSourceConfig> = {
  id: "timer-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Timer",
    label: "timer",
    type: "source",
    description: "Execute the channel at regular intervals",
    io: {
      inputs: [
        {
          id: "passThrough",
          allowedEdge: ["channel"],
          dataType: ["signal"],
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
    config: {
      cronExpression: '*/5 * * * * *'
    },
    dependencies: {
      models: [],
      things: [],
    },
  },
}
