'use client'

import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import { FC } from 'react'
import HttpActionOptions from './HttpAction'
import HttpSourceOptions from './HttpSource'
import MqttActionOptions from './MqttAction'
import TransformProcessingOptions from './TransformProcessing'
import MqttSourceOptions from './MqttSource'
import ThresholdConditionalOptions from './ThresholdConditional'
import TimerSourceOptions from './TimerSource'
import FileStorageActionOptions from './FileStorageAction'
import MessageQueueActionOptions from './MessageQueueAction'
import TimeSeriesDbActionOptions from './TimeSeriesDbAction'
import MessageQueueSourceOptions from './MessageQueueSource'
import NotificationActionOptions from './NotificationAction'
import DocumentDbActionOptions from './DocumentDbAction'
import ExistenceConditionalOptions from './ExistenceConditional'
import StringCompareConditionalOptions from './StringCompareConditional'
import BooleanConditionalOptions from './BooleanConditional'
import HttpResponseActionOptions from './HttpResponseAction'
import MetricActionOptions from './MetricAction'

type Props = {
  spaceId: string
  orgId: string
  channelId: string
}

const NodeOptionsPane: FC<Props> = ({ spaceId, orgId, channelId }) => {
  // Set up store
  const currentNode = useChannelEditorStore((state) => state.currentNode)
  const validationErrors = useChannelEditorStore((state) => state.validationErrors)
  const models = useChannelEditorStore((state) => state.models)
  const things = useChannelEditorStore((state) => state.things)

  const nodeOptions = {
    source: {
      http: <HttpSourceOptions channelId={channelId} />,
      mqtt: <MqttSourceOptions things={things} spaceId={spaceId} orgId={orgId} />,
      messageQueue: <MessageQueueSourceOptions things={things} spaceId={spaceId} orgId={orgId} />,
      timer: <TimerSourceOptions />,
    },
    processing: {
      transform: <TransformProcessingOptions models={models} spaceId={spaceId} orgId={orgId} />,
    },
    conditional: {
      threshold: <ThresholdConditionalOptions models={models} spaceId={spaceId} orgId={orgId} />,
      existence: <ExistenceConditionalOptions models={models} />,
      stringCompare: (
        <StringCompareConditionalOptions models={models} spaceId={spaceId} orgId={orgId} />
      ),
      boolean: <BooleanConditionalOptions models={models} spaceId={spaceId} orgId={orgId} />,
    },
    action: {
      http: <HttpActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      httpResponse: <HttpResponseActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      mqtt: <MqttActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      messageQueue: <MessageQueueActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      fileStorage: <FileStorageActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      timeSeriesDb: <TimeSeriesDbActionOptions things={things} spaceId={spaceId} orgId={orgId} />,
      notification: (
        <NotificationActionOptions
          models={models}
          things={things}
          spaceId={spaceId}
          orgId={orgId}
        />
      ),
      documentDb: (
        <DocumentDbActionOptions models={models} things={things} spaceId={spaceId} orgId={orgId} />
      ),
      log: <></>,
      metric: (
        <MetricActionOptions models={models} things={things} orgId={orgId} spaceId={spaceId} />
      ),
    },
    invalid: <p>The selected node is invalid.</p>,
  }

  return (
    <div className="px-8 py-4 text-sm">
      <h2 className="uppercase text-gray-500 font-bold text-xs mb-2">Node Options</h2>
      {currentNode && (
        <>
          {/* Select the specific options component to render */}
          {nodeOptions[currentNode.metadata.type][currentNode.metadata.label] ??
            nodeOptions.invalid}
        </>
      )}
      {validationErrors &&
        validationErrors.nodeErrors &&
        currentNode &&
        validationErrors.nodeErrors[currentNode.id] &&
        validationErrors.nodeErrors[currentNode.id].length > 0 && (
          <div className="mt-4">
            <h2 className="uppercase text-gray-500 font-bold text-xs mb-2">Node Errors</h2>
            <ul className="list-disc list-inside">
              {validationErrors.nodeErrors[currentNode.id].map((error) => (
                <li key={error}>{error}</li>
              ))}
            </ul>
          </div>
        )}
    </div>
  )
}

export default NodeOptionsPane
