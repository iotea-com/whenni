import { ChannelNode } from "../../../../../node"

export type AwsSNSActionSubnodeConfig = {
  "awsSNS::thing"?: string
  topicARN?: string
  "input::model"?: string
  templateMessage?: string
  defaultValues?: Record<string, any>
}

export const defaultAwsSNSNode: ChannelNode<AwsSNSActionSubnodeConfig> = {
  id: "aws-sns-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "AWS SNS",
    label: "aws-sns",
    type: "action",
    description: "Sends SMS to user",
    io: {
      inputs: [
        {
          id: "notificationInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
      }
      ],
      outputs: []
    },
    config: {
      "awsSNS::thing": "string",
      topicARN: "some-ARN",
      "input::model": "string",
      templateMessage: "select",
      defaultValues: {
        exampleKey: "exampleValue", // Example default value
      }
    },
    dependencies: {
      things: [],
      models: [],
    },
  },
}

