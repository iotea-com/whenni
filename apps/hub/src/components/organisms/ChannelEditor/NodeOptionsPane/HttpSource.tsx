import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode, ChannelNodeConfig, defaultNodes } from '@gruent/libs/engine/nodes/v1'
import { HttpMethod } from '@gruent/libs/engine/nodes/v1/src/action/http'
import { HttpSourceNodeConfig } from '@gruent/libs/engine/nodes/v1/src/source/http'
import { FC, useEffect, useState } from 'react'
import NodeLibrarySectionNode from '../NodeLibraryPane/NodeLibrarySectionNode'

type Props = {
  channelId: string
}

const HttpSourceOptions: FC<Props> = ({ channelId }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<HttpSourceNodeConfig>,
  )

  const [selectedMethod, setSelectedMethod] = useState<HttpMethod>(
    (currentNode.metadata.config.method as HttpMethod) || HttpMethod.GET,
  )
  const [responseTimeout, setResponseTimeout] = useState<number>(
    currentNode.metadata.config.responseTimeout || 10,
  )

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }

    setSelectedMethod((currentConfig.method as HttpMethod) || HttpMethod.GET)
  }, [currentNode])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedMethod) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'httpServer::thing': channelId,
      method: selectedMethod,
      responseTimeout,
    }

    // update dependencies array
    let httpServerThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'httpServer::thing',
    )
    if (httpServerThingDependencyIndex < 0)
      httpServerThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[httpServerThingDependencyIndex] = {
      fieldName: 'httpServer::thing',
      internal: true,
      thingId: channelId,
    }

    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [selectedMethod, responseTimeout, channelId, upsertNode])

  return (
    <>
      <FormFieldSelect
        name="method"
        label="Method"
        options={[
          {
            label: HttpMethod.GET,
            value: HttpMethod.GET,
          },
          {
            label: HttpMethod.POST,
            value: HttpMethod.POST,
          },
          {
            label: HttpMethod.PUT,
            value: HttpMethod.PUT,
          },
          {
            label: HttpMethod.DELETE,
            value: HttpMethod.DELETE,
          },
          {
            label: HttpMethod.HEAD,
            value: HttpMethod.HEAD,
          },
          {
            label: HttpMethod.PATCH,
            value: HttpMethod.PATCH,
          },
        ]}
        value={selectedMethod}
        onChange={(e) => setSelectedMethod(e.target.value as HttpMethod)}
      />
      <FormFieldText
        name="responseTimeout"
        label="Response Timeout"
        inputType="number"
        min="1"
        max="30"
        value={responseTimeout.toString()}
        onChange={(e) => setResponseTimeout(parseInt(e.target.value))}
      />
      <NodeLibrarySectionNode
        label="HTTP Response-action"
        defaultNode={defaultNodes.get('HTTP Response-action') as ChannelNode<ChannelNodeConfig>}
      />
    </>
  )
}

export default HttpSourceOptions
