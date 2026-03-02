'use client'

import { ThingCategoryOptions } from '@iotea/hub/components/modals/CreateThingModal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useSearchParams } from 'next/navigation'
import {
  forwardRef,
  useCallback,
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef,
  useState,
} from 'react'
import MqttClient from '../ThingFormPartials/MqttClient'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import HttpServer from '../ThingFormPartials/HttpServer'
import MqttBroker from '../ThingFormPartials/MqttBroker'
import { handleAddThing, handleHealthcheckThing } from '@iotea/hub/actions/things'
import KafkaCluster from '../ThingFormPartials/KafkaCluster'
import KafkaProducer from '../ThingFormPartials/KafkaProducer'
import NatsServer from '../ThingFormPartials/NatsServer'
import NatsClient from '../ThingFormPartials/NatsClient'
import InfluxDbDatabase from '../ThingFormPartials/InfluxDbDatabase'
import ClickhouseDatabase from '../ThingFormPartials/ClickhouseDb'
import KafkaConsumer from '../ThingFormPartials/KafkaConsumer'
import S3Bucket from '../ThingFormPartials/S3Bucket'
import MinioBucket from '../ThingFormPartials/MinioBucket'
import AwsSESEndpoint from '../ThingFormPartials/AwsSesEndpoint'
import AwsSNSEndpoint from '../ThingFormPartials/AwsSnsEndpoint'
import SendgridClient from '../ThingFormPartials/SendgridClient'
import MongoDbServer from '../ThingFormPartials/MongoDbServer'
import FormFieldTextArea from '@iotea/libs/frontend/components/atoms/FormFieldTextArea'

type Props = {
  orgId: string
  spaceId: string
  secrets: { name: string }[]
  name?: string
  category?: string
  attributes?: Map<string, any>
  hideSubmit?: boolean
  importMode?: boolean
}

type ThingImport = {
  name: string
  category: string
  attributes: Record<string, any>
}

const CreateThingForm = forwardRef(
  (
    {
      orgId,
      spaceId,
      secrets,
      name: initialName,
      category: initialCategory,
      attributes: initialAttributes,
      importMode: initialImportMode = false,
      hideSubmit = false,
    }: Props,
    ref,
  ) => {
    const formRef = useRef<HTMLFormElement>(null)

    const searchParams = useSearchParams()

    const [importMode, setImportMode] = useState<Boolean>(initialImportMode)

    const [name, setName] = useState(initialName || (searchParams.get('name') ?? ''))
    const [category, setCategory] = useState(
      initialCategory || (searchParams.get('category') ?? ''),
    )
    const [currentCategoryFormData, setCurrentCategoryFormData] = useState<
      Map<string, string | number | string[] | number[]>
    >(initialAttributes || new Map())

    const [thingJson, setThingJson] = useState<string>()
    const [thingImportError, setThingImportError] = useState<string>()

    const thingImport = useMemo<ThingImport | null>(() => {
      try {
        if (!thingJson) return null

        const thingImport = JSON.parse(thingJson)

        if (!thingImport.name) {
          setThingImportError('name is required')
          return null
        }
        if (!thingImport.category) {
          setThingImportError('category is required')
          return null
        }
        if (!thingImport.attributes) {
          setThingImportError('attributes are required')
          return null
        }

        setThingImportError(undefined)

        return thingImport as ThingImport
      } catch (_err) {
        setThingImportError('could not parse import object')
        return null
      }
    }, [thingJson])

    useEffect(() => {
      if (thingImport && !thingImportError) {
        setName(thingImport.name)
        setCategory(thingImport.category)
        setCurrentCategoryFormData(new Map(Object.entries(thingImport.attributes)))
      }
    }, [thingImport, thingImportError])

    const handleSubmit = useCallback(async () => {
      if (!formRef.current) return { error: 'Form element not found' }

      const formData = new FormData(formRef.current as HTMLFormElement)
      const name = formData.get('name') as string
      const category = formData.get('category') as string
      const attributes = Object.fromEntries(currentCategoryFormData.entries())

      const { error } = await handleAddThing(spaceId, name, category, attributes)

      if (error) {
        addToast({
          title: 'Could not create the thing',
          body: `The thing could not be created: ${error}`,
          level: 'error',
        })
        return { error: error }

        // TODO: update form hints with invalid fields
      }

      addToast({
        title: 'Successfully added thing to space',
        body: 'Successfully added thing to space.',
        level: 'success',
      })
      return { error: null }
    }, [formRef, currentCategoryFormData, spaceId])

    useImperativeHandle(ref, () => ({
      handleSubmit,
    }))

    const handleHealthcheck = useCallback(async () => {
      const attributes = Object.fromEntries(currentCategoryFormData.entries())
      const { error } = await handleHealthcheckThing(spaceId, category, attributes)

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
    }, [currentCategoryFormData, spaceId, category])

    const secretOptions = useMemo(() => {
      return secrets.map((secret) => ({
        label: secret.name,
        value: secret.name,
      }))
    }, [secrets])

    // If initialAttributes were passed in, create an initialThing object
    const initialThing = useMemo(() => {
      if (!initialAttributes) return undefined

      return {
        attributes: Object.fromEntries(initialAttributes.entries()),
      }
    }, [initialAttributes])

    const currentCategoryFormFields = useMemo(() => {
      switch (category) {
        case 'MQTT_CLIENT':
          return (
            <MqttClient
              orgId={orgId}
              spaceId={spaceId}
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'MQTT_BROKER':
          return <MqttBroker setFormData={setCurrentCategoryFormData} initial={initialThing} />
        case 'HTTP_SERVER':
          return <HttpServer setFormData={setCurrentCategoryFormData} initial={initialThing} />
        case 'KAFKA_CLUSTER':
          return <KafkaCluster setFormData={setCurrentCategoryFormData} initial={initialThing} />
        case 'KAFKA_PRODUCER':
          return (
            <KafkaProducer
              orgId={orgId}
              spaceId={spaceId}
              setFormData={setCurrentCategoryFormData}
              initial={initialThing}
            />
          )
        case 'KAFKA_CONSUMER':
          return (
            <KafkaConsumer
              spaceId={spaceId}
              setFormData={setCurrentCategoryFormData}
              initial={initialThing}
            />
          )
        case 'NATS_SERVER':
          return <NatsServer setFormData={setCurrentCategoryFormData} initial={initialThing} />
        case 'NATS_CLIENT':
          return (
            <NatsClient
              orgId={orgId}
              spaceId={spaceId}
              setFormData={setCurrentCategoryFormData}
              initial={initialThing}
            />
          )
        case 'INFLUXDB_DATABASE':
          return (
            <InfluxDbDatabase
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'CLICKHOUSE_DATABASE':
          return (
            <ClickhouseDatabase
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'S3_BUCKET':
          return (
            <S3Bucket
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'MINIO_BUCKET':
          return (
            <MinioBucket
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'AWS_SNS_ENDPOINT':
          return (
            <AwsSNSEndpoint
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'AWS_SES_ENDPOINT':
          return (
            <AwsSESEndpoint
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'SENDGRID_CLIENT':
          return (
            <SendgridClient
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        case 'MONGODB_SERVER':
          return (
            <MongoDbServer
              setFormData={setCurrentCategoryFormData}
              secretOptions={secretOptions}
              initial={initialThing}
            />
          )
        default:
          return
      }
    }, [category, orgId, spaceId, secretOptions, initialThing])

    const showHealthcheckButton = useMemo(() => {
      switch (category) {
        case 'MQTT_CLIENT':
        case 'KAFKA_PRODUCER':
        case 'KAFKA_CONSUMER':
        case 'NATS_CLIENT':
          return false
        default:
          return true
      }
    }, [category])

    return (
      <>
        {importMode && (
          <section className="mb-8">
            <FormFieldTextArea
              label="Import object"
              name="thingJson"
              value={thingJson}
              onChange={(e) => setThingJson(e.target.value)}
              className="h-40"
            />
            {thingImportError && (
              <p className="text-red-400">Invalid thing import object: {thingImportError}</p>
            )}
          </section>
        )}
        {(!importMode || (importMode && thingImport)) && (
          <section
            className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-700"
            style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
          >
            <form
              ref={formRef}
              onSubmit={(e) => {
                e.preventDefault()
                handleSubmit()
              }}
            >
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
              <FormFieldSelect
                name="category"
                label="Category"
                className="mt-4 mb-2"
                options={ThingCategoryOptions}
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              />
              {currentCategoryFormFields}
              <div className="flex gap-2">
                {!hideSubmit && <Button type="submit" text="Create" />}
                {showHealthcheckButton && (
                  <Button text="Test" variant="secondary" onClick={handleHealthcheck} />
                )}
              </div>
            </form>
          </section>
        )}
        <div className="flex justify-end">
          <Button variant="underline" onClick={() => setImportMode(!importMode)}>
            {importMode ? 'Create mode' : 'Import mode'}
          </Button>
        </div>
      </>
    )
  },
)

CreateThingForm.displayName = 'CreateThingForm'

export default CreateThingForm
