import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import ioteaClient from '@iotea/hub/lib/iotea'
import DeleteThingModal from '@iotea/hub/components/modals/DeleteThingModal'
import PolicySection from './PolicySection'
import { Policy } from '@iotea/libs/iotea-js/src/policies'
import ThingSettingsForm from '@iotea/hub/components/organisms/ThingSettingsForm'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import ExportThingButton from '@iotea/hub/components/molecules/ExportThingButton'
import Container from '@iotea/libs/frontend/components/templates/Container'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'

export const metadata = {
  title: 'Thing details | IOTEA',
}

const ThingDetailsPage = async ({ params }) => {
  const { orgId, spaceId, thingId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: secrets, errors } = await ioteaClient(accessToken).secrets.list(spaceId)

  if (errors && errors.length > 0)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not get the secrets in the space"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )

  if (!secrets)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not get the secrets"
          description={`The secrets could not be retrieved. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )

  const { data: thing, errors: getThingErrors } = await ioteaClient(accessToken).things.get(
    spaceId,
    thingId,
  )

  if (getThingErrors && getThingErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the thing details"
          description={`${getThingErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!thing || !thing.attributes) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the thing details"
          description={`The thing details could not be retrieved. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  const attributes = JSON.parse(thing.attributes.toString())

  const { policy: policyResponse, policyErrors } = await (async () => {
    if (!attributes || !attributes.certificateId) return { policy: null, policyErrors: null }

    const { data: policy, errors: policyErrors } = await ioteaClient(accessToken).policies.get(
      spaceId,
      attributes.certificateId,
    )

    if (policyErrors && policyErrors.length > 0) return { policy: null, policyErrors }
    return { policy, policyErrors: null }
  })()

  if (policyErrors && policyErrors.length > 0)
    return (
      <>
        <Callout
          className="max-w-xl"
          title="Could not get the policy"
          description={`${policyErrors[0]}`}
          variant="error"
        />
      </>
    )

  const policy = (() => {
    if (policyResponse && policyResponse.policy)
      return JSON.parse(policyResponse.policy.toString()) as Policy
  })()

  const revoke = policyResponse?.revoke

  return (
    <>
      <DeleteThingModal orgId={orgId} spaceId={spaceId} thingId={thingId} />
      <Container className="mb-6">
        <h1 className="text-xl">{thing.name}</h1>
        <p className="text-xs">Modify the connection details of your thing.</p>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          {thing && (
            <>
              <h2 className="mb-4">Attributes</h2>
              <ThingSettingsForm orgId={orgId} spaceId={spaceId} thing={thing} secrets={secrets} />
              {thing.thingCategory === 'MQTT_CLIENT' && (
                <div>
                  {policy && (
                    <PolicySection
                      policy={policy}
                      revoke={revoke ?? false}
                      spaceId={spaceId}
                      certificateId={attributes.certificateId}
                      className="my-4"
                    />
                  )}
                </div>
              )}
            </>
          )}
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <FormFieldText label="Thing ID" name="thingId" value={thing.id} disabled />
          <hr className="my-6" />
          <h2 className="mb-2">Export</h2>
          <ExportThingButton
            name={thing.name}
            category={thing.thingCategory}
            attributes={attributes}
          />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-red-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2>Danger Zone</h2>
          <p className="text-xs mb-2">
            These settings may have significant and irreversible changes to your thing.
          </p>
          <Button modalId="deleteThing" text="Delete thing" className="bg-red-500 border-red-500" />
        </section>
      </Container>
    </>
  )
}

export default ThingDetailsPage
