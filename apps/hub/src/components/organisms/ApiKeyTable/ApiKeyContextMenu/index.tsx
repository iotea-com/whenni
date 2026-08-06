'use client'

import { FC } from 'react'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import ContextMenu from '@gruent/hub/components/atoms/ContextMenu'
import { handleDeleteApiKey } from '@gruent/hub/actions/apiKeys'
import { riClipboardLine, riDeleteBin7Line } from '@mwarnerdotme/react-remixicon'
import copyToClipboard from '@gruent/libs/frontend/util/copyToClipboard'

type Props = {
  apiKeyId: string
  orgId: string
  spaceId?: string
  className?: string
}

const ApiKeyContextMenu: FC<Props> = ({ apiKeyId, orgId, spaceId }) => {
  const contextMenuButtons = [
    {
      icon: riClipboardLine,
      title: 'Copy',
      onClick: async () => {
        const { error } = await copyToClipboard(apiKeyId)
        if (error)
          addToast({
            title: 'Could not copy text',
            body: 'Could not copy the API key to your clipboard. This may be due to a browser permission issue.',
          })
        else
          addToast({
            title: 'Copied text',
            body: 'The API key should now be in your clipboard!',
          })
      },
    },
    {
      icon: riDeleteBin7Line,
      title: 'Delete',
      onClick: async () => {
        const { error } = await handleDeleteApiKey(orgId, apiKeyId, spaceId)
        if (error) {
          addToast({
            title: 'Could not delete the API key',
            body: `The API key that you selected could not be deleted: ${error}`,
            level: 'error',
          })
        } else {
          addToast({
            title: 'Successfully deleted API key',
            body: 'Successfully deleted API key from your space.',
            level: 'success',
          })
        }
      },
    },
  ]

  return <ContextMenu contextMenuButtons={contextMenuButtons} />
}

export default ApiKeyContextMenu
