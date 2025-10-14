import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import DeleteSpaceModal from '@iotea/hub/components/modals/DeleteSpaceModal'
import SpaceSettingsForm from '@iotea/hub/components/organisms/SpaceSettingsForm'
import ioteaClient from '@iotea/hub/lib/iotea'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import Container from '@iotea/libs/frontend/components/templates/Container'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Space settings | IOTEA',
}

const SpaceSettingsPage = async ({ params }) => {
  const { orgId, spaceId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: space, errors } = await ioteaClient(accessToken).spaces.get(orgId, spaceId)

  if (errors && errors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the space details"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!space) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the space details"
          description={`The space details could not be retrieved. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <DeleteSpaceModal space={space} />
      <Container className="mb-4">
        <h1 className="text-xl">Space Settings</h1>
        <p>
          <small>Change how your space functions.</small>
        </p>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <SpaceSettingsForm orgId={orgId} space={space} />
          <hr className="my-6" />
          <FormFieldText label="Space Id" name="spaceId" value={space.id} disabled />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-red-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2>Danger Zone</h2>
          <p className="text-xs mb-2">
            These settings may have significant and irreversible changes to your space.
          </p>
          <Button
            id="deleteSpaceButton"
            text="Delete space"
            modalId="deleteSpace"
            className="bg-red-500 border-red-500"
          />
        </section>
      </Container>
    </>
  )
}

export default SpaceSettingsPage
