import { useEffect, useState } from 'react'
import isBrowser from '../util/isBrowser'

export const NotificationEventName = 'toastNotification'

export interface NotificationEvent {
  action: 'add' | 'remove'
}

export type AddNotificationEvent = NotificationEvent & {
  notification: Notification
}

export type RemoveNotificationEvent = NotificationEvent & {
  key: string
}

export type Notification = {
  key: string
  title: string
  body: string
  ttl?: number
  level?: 'info' | 'error' | 'warning' | 'success'
}

export const removeToast = (key: string) => {
  const e = new CustomEvent(NotificationEventName, {
    detail: {
      action: 'remove',
      key,
    },
  })

  window.dispatchEvent(e)
}

export const addToast = ({ title, body, level = 'info', ttl = 5 }: Omit<Notification, 'key'>) => {
  if (!isBrowser) return

  const key = `${title}-${new Date().getTime()}`
  const e = new CustomEvent(NotificationEventName, {
    detail: {
      action: 'add',
      notification: {
        key,
        title,
        body,
        level,
      },
    },
  })

  window.dispatchEvent(e)

  if (ttl > 0)
    setTimeout(() => {
      removeToast(key)
    }, ttl * 1e3)
}

const useToast = () => {
  const [notifications, setNotifications] = useState<Notification[]>([])

  useEffect(() => {
    if (isBrowser) {
      const eventHandler = ((e: CustomEvent) => {
        const notificationEvent = e.detail as NotificationEvent

        if (notificationEvent.action && notificationEvent.action === 'add') {
          const { notification } = notificationEvent as AddNotificationEvent
          setNotifications((current) => {
            return [...current, notification]
          })
        }

        if (notificationEvent.action && notificationEvent.action === 'remove') {
          const { key } = notificationEvent as RemoveNotificationEvent

          setNotifications((current) => {
            const filteredNotifications = current.filter((notification) => {
              if (notification.key === key) return false
              return true
            })

            return filteredNotifications
          })
        }
      }) as EventListener

      window.addEventListener(NotificationEventName, eventHandler)

      return () => {
        window.removeEventListener(NotificationEventName, eventHandler)
      }
    }
  }, [])

  return notifications
}

export default useToast
