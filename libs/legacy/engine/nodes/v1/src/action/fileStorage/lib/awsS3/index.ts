import { ChannelNode } from "../../../../../node"

export type S3ActionSubnodeConfig = {
  "s3bucket::thing"?: string
  key?: string
  keyValues?: string[]
}

export const defaultS3ActionNode: ChannelNode<S3ActionSubnodeConfig> = {
  id: "S3-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "S3",
    label: "s3",
    type: "action",
    description: "Stores data into an S3 bucket",
    io: {
      inputs: [
        {
          id: "fileStorageInput",
          allowedEdge: ["*"],
          dataType: ["bytes", "json", "file"],
      }
      ],
      outputs: []
    },
    config: {
      "s3bucket::thing": "string",
      key: "test-file",
      keyValues: [],
    },
    dependencies: {
      things: [],
      models: [],
    },
  },
}