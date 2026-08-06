'use client'

import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import Modal from '@gruent/libs/frontend/components/organisms/Modal'

const ChannelValidationErrorModal = () => {
  const validationErrors = useChannelEditorStore((state) => state.validationErrors)

  return (
    <Modal id="channelValidationErrors" showAccept={false}>
      <h2>Channel Validation Errors</h2>
      {validationErrors &&
        validationErrors.channelErrors &&
        validationErrors.channelErrors.length > 0 && (
          <ul className="list-disc my-1 pl-4">
            {validationErrors.channelErrors.map((ve, i) => {
              return <li key={i}>{ve}</li>
            })}
          </ul>
        )}
      {(!validationErrors ||
        (validationErrors.channelErrors && validationErrors.channelErrors.length <= 0)) && (
        <p>This channel is valid and ready to publish!</p>
      )}
    </Modal>
  )
}

export default ChannelValidationErrorModal
