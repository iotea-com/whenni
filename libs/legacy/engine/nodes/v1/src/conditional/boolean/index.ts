import { ChannelNode } from "../../../node"

export type BooleanCondition = {
  attributeId: string
  operator: "equal-to"
  value: boolean
}

export type BooleanConditionalNodeConfig = {
  "model::model": string
  conditions: BooleanCondition[]
  logicalOperator: 'AND' | 'OR'
}

export const defaultBooleanConditionalNode: ChannelNode<BooleanConditionalNodeConfig> = {
  id: "boolean-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Boolean",
    label: "boolean",
    type: "conditional",
    description: "Evaluate conditions and route data based on specified rules",
    io: {
      inputs: [
        {
          id: "booleanInput",
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
      conditions: [],
      logicalOperator: "OR"
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
