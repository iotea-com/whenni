'use client'

import { RemixIcon, riCloseFill } from '@mwarnerdotme/react-remixicon'
import type { FC } from 'react'
import { useMemo } from 'react'
import type { Notification } from '../../../hooks/useToast'
import { removeToast } from '../../../hooks/useToast'

type Props = Omit<Notification, 'key' | 'ttl'> & { keyName: string }

const ToastNotification: FC<Props> = ({ keyName, title, body, level }) => {
  const borderColor = useMemo(() => {
    switch (level) {
      case 'success':
        return 'border-success'
      case 'warning':
        return 'border-warning'
      case 'error':
        return 'border-error'
      case 'info':
      default:
        return 'border-info'
    }
  }, [level])

  return (
    <div
      className={`toastNotification relative border-l-4 pl-4 ${borderColor} w-80 bg-gray-50 dark:bg-gray-900 rounded-xs px-4 pt-2 pb-4 shadow-md`}
    >
      <RemixIcon
        icon={riCloseFill}
        className="absolute top-2 right-2 cursor-pointer dark:text-gray-200"
        onClick={() => removeToast(keyName)}
      />
      <span className="font-semibold text-sm dark:text-gray-200">{title}</span>
      <p className="text-xs dark:text-gray-300">{body}</p>
    </div>
  )
}

export default ToastNotification
