import { Thing } from "@prisma/client"
import { ChannelNode } from "../../../node"

export type MqttSourceNodeConfig = {
  "mqttClient::thing": string | Thing
  topic: string
  qos: 0 | 1 | 2
}

export const defaultMqttSourceNode: ChannelNode<MqttSourceNodeConfig> = {
  id: "mqtt-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "MQTT",
    label: "mqtt",
    type: "source",
    description: "Subscribe to messages from an MQTT topic",
    io: {
      inputs: [
        {
          id: "passThrough",
          allowedEdge: ["channel"],
          dataType: ["bytes", "json"],
        }
      ],
      outputs: [
        {
          id: 'mqttOutput',
          allowedEdge: ["processing", "action", "conditional", "custom"],
          dataType: ["bytes", "json"],
        }
      ],
    },
    config: {
      "mqttClient::thing": "",
      topic: "test/topic",
      qos: 0,
    },
    dependencies: {
      models: [],
      things: [
        {
          "internal": false,
          "thingId": "",
          "fieldName": "mqttClient::thing"
        }
      ],
    },
  },
}
