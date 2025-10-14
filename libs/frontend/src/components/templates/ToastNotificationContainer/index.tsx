'use client'

import type { FC } from 'react'
import { useMemo } from 'react'
import useToast from '../../../hooks/useToast'
import ToastNotification from '../../molecules/ToastNotification'

type Props = {
  position?: 'top-right' | 'bottom-right' | 'bottom-left' | 'top-left'
}

const ToastNotificationContainer: FC<Props> = ({ position = 'bottom-right' }) => {
  const notifications = useToast()

  const positionClassNames = useMemo(() => {
    switch (position) {
      case 'top-right':
        return 'top-10 right-10'
      case 'bottom-left':
        return 'bottom-10 left-10'
      case 'top-left':
        return 'top-10 left-10'
      case 'bottom-right':
      default:
        return 'bottom-10 right-10'
    }
  }, [position])

  return (
    <div
      id="toastNotificationContainer"
      className={`fixed flex flex-col gap-2 ${positionClassNames} z-50`}
    >
      {notifications.length > 0 &&
        notifications.map(({ key, title, body, level }) => {
          return (
            <ToastNotification key={key} keyName={key} title={title} body={body} level={level} />
          )
        })}
    </div>
  )
}

export default ToastNotificationContainer
