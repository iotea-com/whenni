'use client'

import { FC } from 'react'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import ContextMenu from '@iotea/hub/components/atoms/ContextMenu'
import {
  handleDeleteChannel,
  handlePublishChannel,
  handleUnpublishChannel,
} from '@iotea/hub/actions/channels'
import { Channel } from '@prisma/client'
import { riCloudLine, riCloudOffLine, riDeleteBin7Line } from '@mwarnerdotme/react-remixicon'

type Props = {
  channel: Channel
  spaceId: string
  className?: string
}

const ChannelContextMenu: FC<Props> = ({ channel, spaceId }) => {
  const publishOption = (() => {
    if (!channel.publishedAt)
      return {
        icon: riCloudLine,
        title: 'Publish',
        onClick: async () => {
          const { error } = await handlePublishChannel(spaceId, channel.id)
          if (error) {
            addToast({
              title: 'Could not publish the channel',
              body: `The channel could not be published: ${error}`,
              level: 'error',
            })
          } else {
            addToast({
              title: 'Publish request sent',
              body: 'The channel should be published soon. This may take a few seconds.',
              level: 'success',
            })
          }
        },
      }

    return {
      icon: riCloudOffLine,
      title: 'Unpublish',
      onClick: async () => {
        const { error } = await handleUnpublishChannel(spaceId, channel.id)
        if (error) {
          addToast({
            title: 'Could not unpublish the channel',
            body: `The channel could not be unpublished: ${error}`,
            level: 'error',
          })
        } else {
          addToast({
            title: 'Unpublish request sent',
            body: 'The channel should be unpublished soon. This may take a few seconds.',
            level: 'success',
          })
        }
      },
    }
  })()

  const contextMenuButtons = [
    publishOption,
    {
      icon: riDeleteBin7Line,
      title: 'Delete',
      onClick: async () => {
        const { error } = await handleDeleteChannel(spaceId, channel.id)
        if (error) {
          addToast({
            title: 'Could not delete the channel',
            body: `The channel could not be deleted: ${error}`,
            level: 'error',
          })
        } else {
          addToast({
            title: 'Successfully deleted channel',
            body: 'Successfully deleted the channel from your space.',
            level: 'success',
          })
        }
      },
    },
  ]

  return <ContextMenu contextMenuButtons={contextMenuButtons} />
}

export default ChannelContextMenu
