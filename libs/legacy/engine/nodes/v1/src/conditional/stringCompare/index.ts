import { ChannelNode } from "../../../node"

export enum StringComparisonOperator {
  "EqualTo" = "equal-to",
  "Contains" = "contains",
  "StartsWith" = "starts-with",
  "EndsWith" = "ends-with",
}

export type StringComparison = {
  attributeId: string
  value: string
  operator: StringComparisonOperator
}

export type StringCompareConditionalNodeConfig = {
  "model::model": string
  comparisons: StringComparison[]
}

export const defaultStringCompareConditionalNode: ChannelNode<StringCompareConditionalNodeConfig> = {
  id: "stringCompare-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "String Compare",
    label: "stringCompare",
    type: "conditional",
    description: "Compare a string attribute in the input data to a value",
    io: {
      inputs: [
        {
          id: "stringCompareInput",
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
      comparisons: []
    },
    dependencies: {
      models: [
        {
          modelId: "",
          attributeName: "model::model",
        },
      ],
      things: [],
    },
  },
}
