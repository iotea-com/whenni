import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import ModelSettingsForm from '@gruent/hub/components/organisms/ModelSettingsForm'
import DeleteModelModal from '@gruent/hub/components/modals/DeleteModelModal'
import gruentClient from '@gruent/hub/lib/gruent'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import ExportModelButton from '@gruent/hub/components/molecules/ExportModelButton'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import Container from '@gruent/libs/frontend/components/templates/Container'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'

export const metadata = {
  title: 'Model details | GRUENT',
}

const ModelDetailsPage = async ({ params }) => {
  const { orgId, spaceId, modelId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: model, errors } = await gruentClient(accessToken).models.get(spaceId, modelId)

  if (errors && errors.length > 0)
    return (
      <Callout
        className="max-w-xl"
        title="Could not get the model details"
        description={`${errors[0]}`}
        variant="error"
      />
    )

  if (!model)
    return (
      <Callout
        className="max-w-xl"
        title="Could not get the model details"
        description={`The model details could not be retrieved. Try logging out and logging back in.`}
        variant="error"
      />
    )
  const { attributes, error: parseError } = (() => {
    try {
      return { attributes: JSON.parse(model.attributes as string) }
    } catch (_err) {
      return { attributes: {}, error: "Could not parse the model's attributes." }
    }
  })()

  return (
    <>
      <DeleteModelModal orgId={orgId} spaceId={spaceId} model={model} />
      <Container className="mb-6">
        <h1 className="text-xl">Model</h1>
        <p className="text-xs">Modify the shape of your data.</p>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2 className="mb-4">Attributes</h2>
          <ModelSettingsForm spaceId={spaceId} initialModel={model} />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <FormFieldText label="Model Id" name="modelId" value={model.id} disabled />
          <hr className="my-6" />
          <h2 className="mb-2">Export</h2>
          {parseError && (
            <Callout
              className="max-w-xl"
              title="Could not export the model"
              description={`${parseError}`}
              variant="warning"
            />
          )}
          {!parseError && <ExportModelButton name={model.name} attributes={attributes} />}
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded-sm py-5 px-10 border border-red-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <h2 className="">Danger Zone</h2>
          <p className="text-xs mb-2">
            These settings may have significant and irreversible changes to your space
          </p>
          <Button variant="primary" className="bg-red-500 border-red-500" modalId="deleteModel">
            <span>Delete model</span>
          </Button>
        </section>
      </Container>
    </>
  )
}

export default ModelDetailsPage
