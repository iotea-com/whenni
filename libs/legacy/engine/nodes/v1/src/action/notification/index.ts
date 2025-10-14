import { ChannelNode } from "../../../node"
import { AwsSNSActionSubnodeConfig } from './lib/awsSNS'
import {  AwsSESActionSubnodeConfig} from './lib/awsSES'

export type NotificationActionNodeConfig = AwsSNSActionSubnodeConfig | AwsSESActionSubnodeConfig 

export const defaultNotificationActionNode: ChannelNode<NotificationActionNodeConfig> = {
  id: "notification-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Notification",
    label: "notification",
    type: "action",
    description: "Send notifications via SMS or email",
    io: {
      inputs: [
        {
          id: "notificationInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
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
