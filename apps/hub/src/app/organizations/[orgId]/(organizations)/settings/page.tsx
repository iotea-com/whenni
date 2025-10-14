import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import DeleteOrganizationModal from '@iotea/hub/components/modals/DeleteOrganizationModal'
import ioteaClient from '@iotea/hub/lib/iotea'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import OrganizationSettingsForm from '@iotea/hub/components/organisms/OrganizationSettingsForm'
import Container from '@iotea/libs/frontend/components/templates/Container'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'

export const metadata = {
  title: 'Organization settings | IOTEA',
}

const OrganizationSettingsPage = async ({ params }) => {
  const { orgId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: organization, errors } = await ioteaClient(accessToken).organizations.get(orgId)

  if (errors && errors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the organization details"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!organization) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Organization not found"
          description={`This organization does not exist. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <DeleteOrganizationModal organization={organization} />
      <Container className="mb-4">
        <h1 className="text-xl">Organization Settings</h1>
        <p>
          <small>Change how your organization functions.</small>
        </p>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <OrganizationSettingsForm organization={organization} />
          <hr className="my-6" />
          <FormFieldText label="Organization Id" name="orgId" value={organization.id} disabled />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-red-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2>Danger Zone</h2>
          <p className="text-xs mb-2">
            These settings may have significant and irreversible changes to your organization.
          </p>
          <Button
            id="deleteOrganizationButton"
            text="Delete organization"
            modalId="deleteOrganization"
            className="bg-red-500 border-red-500"
          />
        </section>
      </Container>
    </>
  )
}

export default OrganizationSettingsPage
