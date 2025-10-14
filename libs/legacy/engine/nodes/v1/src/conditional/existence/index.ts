import { ChannelNode } from "../../../node"

export type ExistenceConditionalNodeConfig = {
  "model::model": string
  attributeIds: string[]
}

export const defaultExistenceConditionalNode: ChannelNode<ExistenceConditionalNodeConfig> = {
  id: "existence-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Existence",
    label: "existence",
    type: "conditional",
    description: "Ensure that a specific field exists in the input data",
    io: {
      inputs: [
        {
          id: "existenceInput",
          allowedEdge: [
            "*"
          ],
          dataType: [
            "bytes",
            "json"
          ]
        }
      ],
      outputs: [
        {
          id: "Fail",
          allowedEdge: [
            "processing",
            "action",
            "conditional",
            "custom"
          ],
          dataType: [
            "bytes",
            "json"
          ]
        },
        {
          id: "Pass",
          allowedEdge: [
            "processing",
            "action",
            "conditional",
            "custom"
          ],
          dataType: [
            "bytes",
            "json"
          ]
        }
      ]
    },
    config: {
      "model::model": "",
      attributeIds: []
    },
    dependencies: {
      models: [
        {
          modelId: "",
          fieldName: "model::model",
        },
      ],
      things: [],
    },
  },
}
