'use client'

import ioteaClient from '@iotea/hub/lib/iotea'
import { Channel } from '@prisma/client'
import { FC, useCallback, useEffect } from 'react'
import useChannelConfig from '@iotea/libs/frontend/hooks/useChannelConfig'
import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import { useDebounce } from '@react-hooks-library/core'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  channel: Channel
  spaceId: string
}

const ChannelOptionsPane: FC<Props> = ({ channel, spaceId }) => {
  const { accessToken } = useAuth()

  // set up store
  const nodes = useChannelEditorStore((state) => state.nodes)
  const edges = useChannelEditorStore((state) => state.edges)
  const notes = useChannelEditorStore((state) => state.notes)
  const runtime = useChannelEditorStore((state) => state.runtime)
  const setRuntime = useChannelEditorStore((state) => state.setRuntime)
  const setValidationErrors = useChannelEditorStore((state) => state.setValidationErrors)

  const channelConfig = useChannelConfig(channel.id, channel.name, runtime, nodes, edges, notes)

  const debouncedChannelConfig = useDebounce(channelConfig, 500) // 500ms delay

  // validate channel on page load and on every config change
  const validateChannel = useCallback(async () => {
    if (!accessToken) return
    if (debouncedChannelConfig.nodes.length === 0) return

    const { data: validationErrors } = await ioteaClient(accessToken).channels.validate(
      spaceId,
      debouncedChannelConfig,
    )

    setValidationErrors(validationErrors) // will be null if there are no validation errors
  }, [debouncedChannelConfig, spaceId, accessToken, setValidationErrors])

  useEffect(() => {
    validateChannel()
  }, [validateChannel])

  // handle change runtime
  const handleChangeRuntimeSize = useCallback(
    (size: 'small' | 'medium' | 'large') => {
      setRuntime({
        ...runtime,
        size,
      })
    },
    [runtime, setRuntime],
  )

  return (
    <div className="flex flex-col px-8 py-4 grow">
      <h2 className="uppercase text-gray-500 font-bold text-xs mb-2">Channel Options</h2>
      {channelConfig && (
        <>
          <FormFieldSelect
            name="runtime"
            label="Runtime Size"
            options={[
              {
                value: 'small',
                label: 'Small',
                description: 'Up to 5 requests per second. Minimum $5/mo.',
              },
              {
                value: 'medium',
                label: 'Medium',
                description: 'Up to 50 requests per second. Minimum $10/mo.',
              },
              {
                value: 'large',
                label: 'Large',
                description: 'Up to 200 requests per second. Minimum $15/mo.',
              },
            ]}
            variant="cards"
            value={channelConfig.runtime?.size || 'small'}
            onChange={(e) =>
              handleChangeRuntimeSize(e.target.value as 'small' | 'medium' | 'large')
            }
          />
        </>
      )}
    </div>
  )
}

export default ChannelOptionsPane
