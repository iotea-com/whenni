import { ChannelNode } from "../../../node"

export type TransformNodeConfig = {
  "inputModel::model": string
  "outputModel::model": string
  mapping: Record<string, string>
}

export const defaultTransformNode: ChannelNode<TransformNodeConfig> = {
  id: "transform-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Data Transform",
    label: "transform",
    type: "processing",
    description: "Convert or modify input data into a specified output format",
    io: {
      inputs: [
        {
          id: "transformInput",
          allowedEdge: ["*"],
          dataType: [
            "bytes"
          ],
        }
      ],
      outputs: [
        {
          id: 'transformOutput',
          allowedEdge: ["*"],
          dataType: [
            "bytes"
          ],
        }
      ]
    },
    config: {
      "inputModel::model": "",
      "outputModel::model": "",
      mapping: {}
    },
    dependencies: {
      things: [],
      models: [
        {
          "modelId": "",
          "fieldName": "inputModel::model"
        },
        {
          "modelId": "",
          "fieldName": "outputModel::model"
        }
      ],
    },
  },
}
