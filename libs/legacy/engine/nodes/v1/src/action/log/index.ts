import { ChannelNode } from "../../../node"

export type LogActionNodeConfig = {}

export const defaultLogActionNode: ChannelNode<LogActionNodeConfig> = {
  id: "log-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Log",
    label: "log",
    type: "action",
    description: "Record data to channel execution logs for monitoring and debugging",
    io: {
      inputs: [
        {
          id: "logInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [],
    },
    config: {},
    dependencies: {
      models: [],
      things: []
    },
  },
}
