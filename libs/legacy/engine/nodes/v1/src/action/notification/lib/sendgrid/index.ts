import { ChannelNode } from "../../../../../node"

export type SendgridActionSubnodeConfig = {
  "sendgridClient::thing"?: string
  sender?: string
  sendAs?: string
  receiver?: string
  subject?: string
  "input::model"?: string
  templateMessage?: string
  defaultValues?: Record<string, any>
}

export const defaultSendgridNode: ChannelNode<SendgridActionSubnodeConfig> = {
  id: "sendgrid-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Sendgrid",
    label: "sendgrid",
    type: "action",
    description: "Sends email messages to email address via Sendgrid.",
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
      "sendgridClient::thing": "string",
      sender: "email@gmail.com",
      receiver: "email2@gmail.com",
      subject: "some-subject",
      "input::model": "string",
      templateMessage: "",
      defaultValues: {}
    },
    dependencies: {
      things: [],
      models: [],
    },
  },
}

