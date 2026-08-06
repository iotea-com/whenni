import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { HttpActionNodeConfig, HttpMethod } from '@gruent/libs/engine/nodes/v1/src/action/http'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

const HttpActionOptions: FC<Props> = ({ things }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<HttpActionNodeConfig>,
  )

  const [selectedHttpServer, setSelectedHttpServer] = useState<Thing>()
  const [selectedHttpServerPath, setSelectedHttpServerPath] = useState<string>()
  const [selectedMethod, setSelectedMethod] = useState<HttpMethod>()

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    const thing = things.find((thing) => {
      if (thing.id === currentConfig['httpServer::thing']) return true
      return false
    })

    setSelectedHttpServer(thing)
    setSelectedHttpServerPath(currentConfig.path)
    setSelectedMethod(currentConfig.method || HttpMethod.GET)
  }, [currentNode, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedHttpServer) return
    if (!selectedHttpServerPath) return
    if (!selectedMethod) return

    // update config
    const updatedNode = structuredClone(currentNode)

    updatedNode.metadata.config = {
      'httpServer::thing': selectedHttpServer.id,
      method: selectedMethod,
      path: selectedHttpServerPath,
    }

    // update dependencies array
    let httpServerThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'httpServer::thing',
    )
    if (httpServerThingDependencyIndex < 0)
      httpServerThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[httpServerThingDependencyIndex] = {
      fieldName: 'httpServer::thing',
      internal: false,
      thingId: selectedHttpServer.id,
    }

    // apply updates to the nodes map
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [selectedHttpServer, selectedHttpServerPath, selectedMethod, upsertNode])

  const httpServers = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => {
      if (thing.thingCategory == 'HTTP_SERVER') return true
      return false
    })
  }, [things])

  const httpServerOptions = useMemo(() => {
    return httpServers.map((httpServer) => ({
      label: httpServer.name,
      value: httpServer.id,
    }))
  }, [httpServers])

  const selectedHttpServerPathOptions = useMemo(() => {
    if (!selectedHttpServer) return []
    if (!selectedHttpServer.attributes) return []
    const attributes = JSON.parse(selectedHttpServer.attributes as string)
    if (!attributes.paths) return []

    return attributes.paths.map((path) => ({
      label: path,
      value: path,
    }))
  }, [selectedHttpServer])

  if (!httpServers || httpServers.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add an HTTP server
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="httpServer::thing"
          label="HTTP Server"
          className="grow"
          options={httpServerOptions}
          value={selectedHttpServer?.id ?? '__GRUENT_IGNORE__'}
          onChange={(e) =>
            setSelectedHttpServer(
              httpServers.find((s) => {
                if (s.id === e.target.value) return true
                return false
              }),
            )
          }
        />
        <Button className="my-3" onClick={() => openModal('createThing')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      <FormFieldSelect
        name="path"
        label="Path"
        disabled={!selectedHttpServer}
        options={selectedHttpServerPathOptions}
        optional
        value={selectedHttpServerPath}
        onChange={(e) => setSelectedHttpServerPath(e.target.value)}
      />
      <FormFieldSelect
        name="method"
        label="Method"
        options={[
          {
            label: 'GET',
            value: 'GET',
          },
          {
            label: 'POST',
            value: 'POST',
          },
          {
            label: 'PUT',
            value: 'PUT',
          },
          {
            label: 'DELETE',
            value: 'DELETE',
          },
          {
            label: 'HEAD',
            value: 'HEAD',
          },
          {
            label: 'PATCH',
            value: 'PATCH',
          },
        ]}
        value={selectedMethod}
        onChange={(e) => setSelectedMethod(e.target.value as HttpMethod)}
      />
    </>
  )
}

export default HttpActionOptions
