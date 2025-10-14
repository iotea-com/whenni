import { ChannelNode } from "../../../node"
import { InfluxDBNodeConfig } from './lib/influxdb'

export type TimeSeriesDbNodeConfig = InfluxDBNodeConfig

export const defaultTimeSeriesDbActionNode: ChannelNode<TimeSeriesDbNodeConfig> = {
  id: "time-series-beta-dev",
  coordinates: {
    x: 0,
    y: 0,
  },
  metadata: {
    version: "beta-dev",
    name: "Time Series Database",
    label: "timeSeriesDb",
    type: "action",
    description: "Query and store data points in a time series database",
    io: {
      inputs: [
        {
          id: "timeSeriesDbInput",
          allowedEdge: ["*"],
          dataType: ["bytes", "json"],
        }
      ],
      outputs: [],
    },
    config: {},
    dependencies: {
      things: [],
      models: [],
    },
  },
}
