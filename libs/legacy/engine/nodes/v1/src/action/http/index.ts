import { Thing } from "@prisma/client"
import { ChannelNode } from "../../../node"

export enum HttpMethod {
  'GET' = 'GET',
  'POST' = 'POST',
  'PUT' = 'PUT',
  'DELETE' = 'DELETE',
  'PATCH' = 'PATCH',
  'HEAD' = 'HEAD',
}

export type HttpActionNodeConfig = {
  "httpServer::thing": string | Thing
  method: HttpMethod
  path: string
  headers?: {
    [key: string]: string
  }
  body?: {
    [key: string]: any
  }
}

export const defaultHttpActionNode: ChannelNode<HttpActionNodeConfig> = {
  id: "http-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "HTTP",
    label: "http",
    type: "action",
    description: "Send an HTTP request",
    io: {
      inputs: [
        {
          id: "httpInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [
        {
          id: 'httpOutput',
          allowedEdge: ["processing", "action", "conditional", "custom"],
          dataType: ["signal"],
        }
      ],
    },
    config: {
      "httpServer::thing": '__IOTEA_IGNORE__',
      method: HttpMethod.GET,
      path: '__IOTEA_IGNORE__',
    },
    dependencies: {
      models: [],
      things: [
        {
          "internal": false,
          "thingId": "",
          "fieldName": "httpServer::thing"
        }
      ]
    },
  },
}
