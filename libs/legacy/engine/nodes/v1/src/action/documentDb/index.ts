import { ChannelNode } from "../../../node"
import { MongoDbActionSubnodeConfig } from './lib/mongodb'

export type DocumentDbActionNodeConfig = MongoDbActionSubnodeConfig

export const defaultDocumentDbActionNode: ChannelNode<DocumentDbActionNodeConfig> = {
  id: "documentDb-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Document Database",
    label: "documentDb",
    type: "action",
    description: "Query and store data in a document database",
    io: {
      inputs: [
        {
          id: "documentDbInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [
        {
          id: "documentDbOutput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ]
    },
    config: {},
    dependencies: {
      things: [],
      models: [],
    },
  },
}
