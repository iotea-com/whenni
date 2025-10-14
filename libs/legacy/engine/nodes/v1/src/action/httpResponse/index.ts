import { ChannelNode } from "../../../node"

export type HttpResponseActionNodeConfig = {
  responseCode: number
  headers?: Record<string, string>
}

export const defaultHttpResponseActionNode: ChannelNode<HttpResponseActionNodeConfig> = {
  id: "http-response-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "HTTP Response",
    label: "httpResponse",
    type: "action",
    description: "Respond to the HTTP source node's request",
    io: {
      inputs: [
        {
          id: "httpResponseInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [
        {
          id: 'httpResponseOutput',
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
    },
    config: {
      responseCode: 200,
      headers: {
        "Content-Type": "application/json",
      },
    },
    dependencies: {
      models: [],
      things: [],
    },
  },
}
