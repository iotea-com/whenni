import { ChannelNode } from "../../../node"
import { Thing, Model } from "@prisma/client"

export type MetricActionNodeConfig = {
  model: string | Model
  "clickhouseDatabase::thing"?: string | Thing
  attribute?: string
  metadata?: string
}

export const defaultMetricActionNode: ChannelNode<MetricActionNodeConfig> = {
  id: "metric-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Metric",
    label: "metric",
    type: "action",
    description: "Save a data point for viewing in the analytics dashboard",
    io: {
      inputs: [
        {
          id: "metricInput",
          allowedEdge: ["*"],
          dataType: ["bytes"],
        }
      ],
      outputs: [],
    },
    config: {
      model: '__IOTEA_IGNORE__',
      "clickhouseDatabase::thing": '__IOTEA_IGNORE__',
      attribute: '__IOTEA_IGNORE__',
      metadata: '__IOTEA_IGNORE__'
    },
    dependencies: {
      models: [
        {
          modelId: '__IOTEA_IGNORE__',
          fieldName: 'model::model'
        }
      ],
      things: [],
    },
  },
}
