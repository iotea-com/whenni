import { ChannelNode } from "../../../../../node"

export type MinioActionSubnodeConfig = {
  "minioBucket::thing"?: string
  key?: string
  keyValues?: string[]
}

export const defaultMinioActionNode: ChannelNode<MinioActionSubnodeConfig> = {
  id: "Minio-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Minio",
    label: "minio",
    type: "action",
    description: "Stores data into a Minio bucket",
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
      "minioBucket::thing": "string",
      key: "test-file",
      keyValues: [],
    },
    dependencies: {
      things: [],
      models: [],
    },
  },
}