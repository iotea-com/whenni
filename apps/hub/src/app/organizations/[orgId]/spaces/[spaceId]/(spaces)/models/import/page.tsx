import ModelImportForm from '@iotea/hub/components/organisms/ModelImportForm'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import ioteaClient from '@iotea/hub/lib/iotea'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Import a model | IOTEA',
}

const ImportModelPage = async ({ params }) => {
  const { orgId, spaceId } = await params

  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to create the channel.',
    }

  const { data: secrets, errors } = await ioteaClient(accessToken).secrets.list(spaceId)

  if (errors && errors.length > 0)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not get the secrets in this space"
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
          title="Could not get the secrets in this space"
          description={`The secrets could not be listed. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )

  return (
    <>
      <h1>Import Model</h1>
      <p className="text-gray-500 text-sm">
        Paste an import object (from a model details page) to create a new model.
      </p>
      <ModelImportForm orgId={orgId} spaceId={spaceId} />
    </>
  )
}

export default ImportModelPage
