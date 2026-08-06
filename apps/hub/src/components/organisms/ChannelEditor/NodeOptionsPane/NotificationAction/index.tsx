import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { NotificationActionNodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/notification'
import { Thing, Model } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import AwsSESActionOptions from './AwsSESActionOptions'
import AwsSNSActionOptions from './AwsSNSActionOptions'
import SendgridActionOptions from './SendgridActionOptions'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  models: Model[]
}

export enum Notification {
  AWS_SES_ENDPOINT = 'AWS_SES_ENDPOINT',
  AWS_SNS_ENDPOINT = 'AWS_SNS_ENDPOINT',
  SENDGRID_CLIENT = 'SENDGRID_CLIENT',
}

const NotificationActionOptions: FC<Props> = ({ things, orgId, spaceId, models }) => {
  // Accessing the current node in the editor store
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<NotificationActionNodeConfig>,
  )

  // State to hold the selected notification thing (AWS SES or SNS) and selected Model
  const [selectedNotificationThing, setSelectedNotificationThing] = useState<Thing>()

  // Load initial config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things || !models) return

    // Check if an AWS SES thing is already selected in the config
    if (Object.keys(currentConfig).includes('awsSES::thing')) {
      const awsSES = things.find((s) => {
        if (s.id === currentConfig['awsSES::thing']) return true
        return false
      })
      if (awsSES) setSelectedNotificationThing(awsSES)
      return
    }

    // Check if an AWS SNS thing is already selected in the config
    if (Object.keys(currentConfig).includes('awsSNS::thing')) {
      const awsSNS = things.find((s) => {
        if (s.id === currentConfig['awsSNS::thing']) return true
        return false
      })
      if (awsSNS) setSelectedNotificationThing(awsSNS)
      return
    }

    // Check if a Sendgrid client thing is already selected in the config
    if (Object.keys(currentConfig).includes('sendgridClient::thing')) {
      const sendgridClient = things.find((s) => {
        if (s.id === currentConfig['sendgridClient::thing']) return true
        return false
      })
      if (sendgridClient) setSelectedNotificationThing(sendgridClient)
      return
    }
  }, [currentNode, things, models])

  // Determine the notification type based on the selected thing
  const selectedNotificationType = useMemo(() => {
    if (!selectedNotificationThing) return

    if (selectedNotificationThing.thingCategory === 'AWS_SES_ENDPOINT')
      return Notification.AWS_SES_ENDPOINT

    if (selectedNotificationThing.thingCategory === 'AWS_SNS_ENDPOINT')
      return Notification.AWS_SNS_ENDPOINT

    if (selectedNotificationThing.thingCategory === 'SENDGRID_CLIENT')
      return Notification.SENDGRID_CLIENT
  }, [selectedNotificationThing])

  // Generate a filtered list of available notification types (SES or SNS)
  const notificationThings = useMemo(() => {
    if (!things) return []

    const awsSES = things.filter((thing) => {
      if (thing.thingCategory == 'AWS_SES_ENDPOINT') return true
      return false
    })

    const awsSNS = things.filter((thing) => {
      if (thing.thingCategory == 'AWS_SNS_ENDPOINT') return true
      return false
    })

    const sendgridClients = things.filter((thing) => {
      if (thing.thingCategory == 'SENDGRID_CLIENT') return true
      return false
    })

    return [...awsSES, ...awsSNS, ...sendgridClients]
  }, [things])

  // Generate options for the NotificationType dropdown
  const notificationOptions = useMemo(() => {
    return notificationThings.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [notificationThings])

  // Conditionally rendering subnode options based on the selected type (SES/SNS)
  const selectedSubnodeOptions = useMemo(() => {
    if (!selectedNotificationThing) return

    switch (selectedNotificationType) {
      case Notification.AWS_SES_ENDPOINT:
        return (
          <AwsSESActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedNotificationThing={selectedNotificationThing}
            models={models}
          />
        )
      case Notification.AWS_SNS_ENDPOINT:
        return (
          <AwsSNSActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedNotificationThing={selectedNotificationThing}
            models={models}
          />
        )
      case Notification.SENDGRID_CLIENT:
        return (
          <SendgridActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedNotificationThing={selectedNotificationThing}
            models={models}
          />
        )
      default:
        return
    }
  }, [models, selectedNotificationType, selectedNotificationThing, things, orgId, spaceId])

  if (!notificationThings || notificationThings.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a notification client
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="NotificationThing"
          label="Notification Thing"
          className="grow"
          options={notificationOptions}
          value={selectedNotificationThing?.id ?? '__GRUENT_IGNORE__'}
          onChange={(e) =>
            setSelectedNotificationThing(
              notificationThings.find((s) => {
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
      {selectedSubnodeOptions}
    </>
  )
}

export default NotificationActionOptions
