import { ChannelNode } from "../../../node"

export type HttpSourceNodeConfig = {
  method: string
  responseTimeout: number
  "httpServer::thing": string
}

export const defaultHttpSourceNode: ChannelNode<HttpSourceNodeConfig> = {
  id: "http-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "HTTP",
    label: "http",
    type: "source",
    description: "Listen for incoming HTTP requests",
    io: {
      inputs: [
        {
          id: "passThrough",
          allowedEdge: ["channel"],
          dataType: ["bytes", "json", "file"],
        }
      ],
      outputs: [
        {
          id: 'httpOutput',
          allowedEdge: ["*"],
          dataType: ["bytes", "json", "file"],
        }
      ],
    },
    config: {
      "httpServer::thing": "",
      method: "GET",
      responseTimeout: 10,
    },
    dependencies: {
      models: [],
      things: [
        {
          "internal": true,
          "thingId": "",
          "fieldName": "httpServer::thing"
        }
      ]
    },
  },
}
