import { ChannelNode } from "../../../node"
import { S3ActionSubnodeConfig } from './lib/awsS3'
import { MinioActionSubnodeConfig } from './lib/minio'


export type FileStorageActionNodeConfig = S3ActionSubnodeConfig  | MinioActionSubnodeConfig

export const defaultFileStorageActionNode: ChannelNode<FileStorageActionNodeConfig> = {
  id: "file-storage-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "File Storage",
    label: "fileStorage",
    type: "action",
    description: "Save data as files in static file storage",
    io: {
      inputs: [
        {
          id: "fileStorageInput",
          allowedEdge: ["*"],
          dataType: ["bytes", "json", "file"],
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
