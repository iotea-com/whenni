import { Thing } from "@prisma/client"
import { ChannelNode } from "../../../node"

export type MqttActionNodeConfig = {
  "mqttClient::thing": string | Thing
  topic: string
  qos: 0 | 1 | 2
}

export const defaultMqttActionNode: ChannelNode<MqttActionNodeConfig> = {
  id: "mqtt-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "MQTT",
    label: "mqtt",
    type: "action",
    description: "Publish messages to an MQTT topic",
    io: {
      inputs: [
        {
          id: "mqttInput",
          allowedEdge: ["*"],
          dataType: [
            "bytes"
          ],
        }
      ],
      outputs: [],
    },
    config: {
      "mqttClient::thing": "id",
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
