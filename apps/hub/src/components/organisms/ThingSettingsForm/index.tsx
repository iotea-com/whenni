'use client'

import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { FC, FormEvent, useCallback, useMemo, useRef, useState } from 'react'
import MqttClient from '@gruent/hub/components/organisms/ThingFormPartials/MqttClient'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import HttpServer from '@gruent/hub/components/organisms/ThingFormPartials/HttpServer'
import MqttBroker from '@gruent/hub/components/organisms/ThingFormPartials/MqttBroker'
import { handleHealthcheckThing, handleUpdateThing } from '@gruent/hub/actions/things'
import KafkaCluster from '@gruent/hub/components/organisms/ThingFormPartials/KafkaCluster'
import KafkaProducer from '@gruent/hub/components/organisms/ThingFormPartials/KafkaProducer'
import NatsServer from '@gruent/hub/components/organisms/ThingFormPartials/NatsServer'
import NatsClient from '@gruent/hub/components/organisms/ThingFormPartials/NatsClient'
import InfluxDbDatabase from '@gruent/hub/components/organisms/ThingFormPartials/InfluxDbDatabase'
import ClickhouseDatabase from '@gruent/hub/components/organisms/ThingFormPartials/ClickhouseDb'
import KafkaConsumer from '@gruent/hub/components/organisms/ThingFormPartials/KafkaConsumer'
import S3Bucket from '@gruent/hub/components/organisms/ThingFormPartials/S3Bucket'
import MinioBucket from '@gruent/hub/components/organisms/ThingFormPartials/MinioBucket'
import AwsSESEndpoint from '@gruent/hub/components/organisms/ThingFormPartials/AwsSesEndpoint'
import AwsSNSEndpoint from '@gruent/hub/components/organisms/ThingFormPartials/AwsSnsEndpoint'
import SendgridClient from '@gruent/hub/components/organisms/ThingFormPartials/SendgridClient'
import MonogoDbServer from '@gruent/hub/components/organisms/ThingFormPartials/MongoDbServer'
import { Thing } from '@prisma/client'

type Props = {
  orgId: string
  spaceId: string
  thing: Thing
  className?: string
  secrets: { name: string }[]
}

const ThingSettingsForm: FC<Props> = ({ orgId, spaceId, thing, className, secrets }) => {
  const formRef = useRef<HTMLFormElement>(null)

  const [name, setName] = useState(thing.name)
  const [currentCategoryFormData, setCurrentCategoryFormData] = useState<
    Map<string, string | number | string[] | number[]>
  >(new Map())

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const formData = new FormData(formRef.current as HTMLFormElement)
    const name = formData.get('name') as string
    const attributes = Object.fromEntries(currentCategoryFormData.entries())

    const { error } = await handleUpdateThing(spaceId, thing.id, {
      name,
      attributes,
    })

    if (error) {
      addToast({
        title: 'Could not update the thing',
        body: `The thing could not be updated: ${error}`,
        level: 'error',
      })
      return false

      // TODO: update form hints with invalid fields
    }

    addToast({
      title: 'Successfully updated thing',
      body: 'The changes should be reflected immediately.',
      level: 'success',
    })
    return true
  }

  const handleHealthcheck = useCallback(async () => {
    const attributes = Object.fromEntries(currentCategoryFormData.entries())
    const { error } = await handleHealthcheckThing(spaceId, thing.thingCategory, attributes)

    if (error) {
      addToast({
        title: 'Healthcheck failed',
        body: error,
        level: 'warning',
      })
      return
    }

    addToast({
      title: 'Healthcheck successful',
      body: 'A test connection was established. Happy connecting!',
      level: 'success',
    })
  }, [currentCategoryFormData, spaceId, thing.thingCategory])

  const secretOptions = useMemo(() => {
    return secrets.map((secret) => ({
      label: secret.name,
      value: secret.name,
    }))
  }, [secrets])

  const currentCategoryFormFields = useMemo(() => {
    switch (thing.thingCategory) {
      case 'MQTT_CLIENT':
        return (
          <MqttClient
            orgId={orgId}
            spaceId={spaceId}
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'MQTT_BROKER':
        return <MqttBroker setFormData={setCurrentCategoryFormData} initial={thing} />
      case 'HTTP_SERVER':
        return <HttpServer setFormData={setCurrentCategoryFormData} initial={thing} />
      case 'KAFKA_CLUSTER':
        return <KafkaCluster setFormData={setCurrentCategoryFormData} initial={thing} />
      case 'KAFKA_PRODUCER':
        return (
          <KafkaProducer
            orgId={orgId}
            spaceId={spaceId}
            setFormData={setCurrentCategoryFormData}
            initial={thing}
          />
        )
      case 'KAFKA_CONSUMER':
        return (
          <KafkaConsumer
            spaceId={spaceId}
            setFormData={setCurrentCategoryFormData}
            initial={thing}
          />
        )
      case 'NATS_SERVER':
        return <NatsServer setFormData={setCurrentCategoryFormData} initial={thing} />
      case 'NATS_CLIENT':
        return (
          <NatsClient
            orgId={orgId}
            spaceId={spaceId}
            setFormData={setCurrentCategoryFormData}
            initial={thing}
          />
        )
      case 'INFLUXDB_DATABASE':
        return (
          <InfluxDbDatabase
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'CLICKHOUSE_DATABASE':
        return (
          <ClickhouseDatabase
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'S3_BUCKET':
        return (
          <S3Bucket
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'MINIO_BUCKET':
        return (
          <MinioBucket
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'AWS_SNS_ENDPOINT':
        return (
          <AwsSNSEndpoint
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'AWS_SES_ENDPOINT':
        return (
          <AwsSESEndpoint
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'SENDGRID_CLIENT':
        return (
          <SendgridClient
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      case 'MONGODB_SERVER':
        return (
          <MonogoDbServer
            setFormData={setCurrentCategoryFormData}
            initial={thing}
            secretOptions={secretOptions}
          />
        )
      default:
        return
    }
  }, [thing, spaceId, orgId, secretOptions])

  const showHealthcheckButton = useMemo(() => {
    switch (thing.thingCategory) {
      case 'MQTT_CLIENT':
      case 'KAFKA_PRODUCER':
      case 'KAFKA_CONSUMER':
      case 'NATS_CLIENT':
        return false
      default:
        return true
    }
  }, [thing.thingCategory])

  return (
    <div className={className}>
      <form ref={formRef} onSubmit={(e) => handleSubmit(e)}>
        <FormFieldText
          name="name"
          label="Name"
          placeholder="Name of the thing"
          inputType="text"
          className="mt-4 mb-2"
          autocomplete="off"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <FormFieldText
          name="category"
          label="Category"
          className="mt-4 mb-2"
          disabled
          value={thing.thingCategory}
        />
        {currentCategoryFormFields}
        <div className="flex gap-2">
          <Button type="submit" text="Update" />
          {showHealthcheckButton && (
            <Button text="Test" variant="secondary" onClick={handleHealthcheck} />
          )}
        </div>
      </form>
    </div>
  )
}

export default ThingSettingsForm
