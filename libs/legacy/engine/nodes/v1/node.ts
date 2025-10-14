import { MqttSourceNodeConfig, defaultMqttSourceNode } from './src/source/mqtt'
import { MqttActionNodeConfig, defaultMqttActionNode } from './src/action/mqtt'
import { TransformNodeConfig, defaultTransformNode } from './src/processing/transform'
import { HttpActionNodeConfig, defaultHttpActionNode } from './src/action/http'
import { HttpSourceNodeConfig, defaultHttpSourceNode } from './src/source/http'
import { ThresholdConditionalNodeConfig, defaultThresholdConditionalNode } from './src/conditional/threshold'
import { TimerSourceConfig, defaultTimerSourceNode } from './src/source/timer'
import { FileStorageActionNodeConfig, defaultFileStorageActionNode } from './src/action/fileStorage'
import { NotificationActionNodeConfig, defaultNotificationActionNode } from './src/action/notification'
import { DocumentDbActionNodeConfig, defaultDocumentDbActionNode } from './src/action/documentDb'
import { MessageQueueActionNodeConfig, defaultMessageQueueActionNode } from './src/action/messageQueue'
import { TimeSeriesDbNodeConfig, defaultTimeSeriesDbActionNode } from './src/action/timeSeriesDb'
import { MessageQueueSourceNodeConfig, defaultMessageQueueSourceNode } from './src/source/messageQueue'
import { LogActionNodeConfig, defaultLogActionNode } from './src/action/log'
import { defaultExistenceConditionalNode } from './src/conditional/existence'
import { defaultStringCompareConditionalNode } from './src/conditional/stringCompare'
import { BooleanConditionalNodeConfig, defaultBooleanConditionalNode } from './src/conditional/boolean'
import { defaultHttpResponseActionNode } from './src/action/httpResponse'
import { MetricActionNodeConfig, defaultMetricActionNode } from './src/action/metric'

export type ChannelNodeId = string
export type ChannelNodeType = "source" | "processing" | "action" | "conditional" | "custom"
export type ChannelNodeConfig =
  MqttSourceNodeConfig |
  MqttActionNodeConfig |
  TransformNodeConfig |
  HttpActionNodeConfig |
  HttpSourceNodeConfig |
  ThresholdConditionalNodeConfig |
  TimerSourceConfig |
  FileStorageActionNodeConfig |
  MessageQueueActionNodeConfig |
  TimeSeriesDbNodeConfig |
  MessageQueueSourceNodeConfig |
  NotificationActionNodeConfig |
  DocumentDbActionNodeConfig |
  LogActionNodeConfig |
  BooleanConditionalNodeConfig |
  MetricActionNodeConfig


export type ChannelNodeIOAllowedEdge = ChannelNodeType | "channel" | "*"
export type ChannelNodeIODataType = "bytes" | "json" | "file" | "signal" | "boolean" | "integer" | "float" | "string" | "map" | "pass"

export type ChannelNodeIO = {
  id: string
  allowedEdge: ChannelNodeIOAllowedEdge[]
  dataType: ChannelNodeIODataType[]
}

export type ChannelNode<C = ChannelNodeConfig> = {
  id: string
  coordinates: {
    x: number
    y: number
  }
  metadata: {
    version: string
    name: string
    label: string
    type: ChannelNodeType
    description: string
    io: {
      inputs: ChannelNodeIO[]
      outputs: ChannelNodeIO[]
    }
    config: C
    dependencies: {
      models: any[]
      things: any[]
    }
  }
}

export const defaultNodes = new Map<string, ChannelNode>([
  ['MQTT-source', defaultMqttSourceNode],
  ['HTTP-source', defaultHttpSourceNode],
  ['Timer-source', defaultTimerSourceNode],
  ['Data Transform-processing', defaultTransformNode],
  ['HTTP-action', defaultHttpActionNode],
  ['HTTP Response-action', defaultHttpResponseActionNode],
  ['MQTT-action', defaultMqttActionNode],
  ['Threshold-conditional', defaultThresholdConditionalNode],
  ['Existence-conditional', defaultExistenceConditionalNode],
  ['String Compare-conditional', defaultStringCompareConditionalNode],
  ['File Storage-action', defaultFileStorageActionNode],
  ['Message Queue-action', defaultMessageQueueActionNode],
  ['Message Queue-source', defaultMessageQueueSourceNode],
  ['Time Series Database-action', defaultTimeSeriesDbActionNode],
  ['Notification-action', defaultNotificationActionNode],
  ['Document Database-action', defaultDocumentDbActionNode],
  ['Log-action', defaultLogActionNode],
  ['Boolean-conditional', defaultBooleanConditionalNode],
  ['Metric-action', defaultMetricActionNode],
])