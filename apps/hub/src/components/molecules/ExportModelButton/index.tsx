'use client'

import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import copyToClipboard from '@iotea/libs/frontend/util/copyToClipboard'
import { FC, useCallback } from 'react'

type Props = {
  name: string
  attributes: ModelAttributes
}

const ExportModelButton: FC<Props> = ({ name, attributes }) => {
  const handleExport = useCallback(async () => {
    const exportData = {
      name,
      attributes,
    }

    const { error } = await copyToClipboard(JSON.stringify(exportData))

    if (error) {
      addToast({
        title: 'Could not copy to clipboard',
        body: 'This may be due to browser permissions. Please allow copying to clipboard and try again.',
        level: 'error',
      })
    } else {
      addToast({
        title: 'Copied model to clipboard',
        body: 'The model has been copied to your clipboard. You can now paste it into the import model form.',
        level: 'success',
      })
    }
  }, [name, attributes])

  return <Button variant="primary" onClick={handleExport} text="Export model" />
}

export default ExportModelButton
