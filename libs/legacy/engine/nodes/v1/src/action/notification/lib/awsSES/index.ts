import { ChannelNode } from "../../../../../node"

export type AwsSESActionSubnodeConfig = {
  "awsSES::thing"?: string
  sender?: string
  receiver?: string
  subject?: string
  "input::model"?: string
  templateMessage?: string
  defaultValues?: Record<string, any>
}

export const defaultAwsSESNode: ChannelNode<AwsSESActionSubnodeConfig> = {
  id: "aws-ses-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "AWS SES",
    label: "aws-ses",
    type: "action",
    description: "Sends email messages to email address",
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
      "awsSES::thing": "string",
      sender: "email@gmail.com",
      receiver: "email2@gmail.com",
      subject: "some-subject",
      "input::model": "string",
      templateMessage: "select",
      defaultValues: {
        "exampleKey": "exampleValue", // Example default value
      }
    },
    dependencies: {
      things: [],
      models: [],
    },
  },
}

