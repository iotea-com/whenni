import { ChannelNode } from "../../../node"

export enum Operator {
  "GreaterThan" = "greater-than",
  "LessThan" = "less-than",
  "GreaterThanEqualTo" = "greater-than-equal-to",
  "LessThanEqualTo" = "less-than-equal-to",
  "EqualTo" = "equal-to",
  "NotEqualTo" = "not-equal-to",
}

export type Condition = {
  attributeId: string
  operator: Operator
  value: number
}

export type ThresholdConditionalNodeConfig = {
  "model::model": string
  conditions: Condition[]
  logicalOperator: 'AND' | 'OR'
}

export const defaultThresholdConditionalNode: ChannelNode<ThresholdConditionalNodeConfig> = {
  id: "threshold-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Threshold",
    label: "threshold",
    type: "conditional",
    description: "Evaluate conditions and route data based on specified rules",
    io: {
      inputs: [
        {
          id: "thresholdInput",
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
