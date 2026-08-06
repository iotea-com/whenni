import gruentClient from '@gruent/hub/lib/gruent'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import CreateSecretModal from '@gruent/hub/components/modals/CreateSecretModal'
import SecretsTable from '@gruent/hub/components/organisms/SecretsTable'
import { riAddCircleLine, riLockLine } from '@mwarnerdotme/react-remixicon'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
import Container from '@gruent/libs/frontend/components/templates/Container'
import ListTablePlaceholder from '@gruent/hub/components/molecules/ListTablePlaceholder'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'

export const metadata = {
  title: 'Secrets | GRUENT',
}

const SecretsPage = async ({ params }) => {
  const { spaceId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: secrets, errors } = await gruentClient(accessToken).secrets.list(spaceId)

  if (errors && errors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the secrets for this space"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!secrets) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the secrets for this space"
          description={`The secrets could not be listed. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <CreateSecretModal spaceId={spaceId} />
      <Container>
        <section
          className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          {Number(secrets.length) > 0 && (
            <>
              <div className="absolute top-8 right-10 mb-1 flex items-end gap-2">
                <div className="grow" />
                <Button variant="transparent" className="text-sm" modalId="addSecret">
                  <RemixIcon className="mr-1" icon={riAddCircleLine} />
                  <span>Create</span>
                </Button>
              </div>
              <SecretsTable
                secrets={secrets}
                spaceId={spaceId}
                totalPages={Math.max(1, Math.round(secrets.length / 10))}
                totalResults={secrets.length}
              />
            </>
          )}
          {Number(secrets.length) <= 0 && (
            <ListTablePlaceholder
              title="Secrets"
              description="Keys, passwords, tokens - anything that you need to stay private."
            >
              <RemixIcon icon={riLockLine} className="text-gray-600 mb-2" size="3x" />
              <p className="text-gray-600 text-center mb-2">No secrets... yet!</p>
              <Button variant="primary" className="text-sm" modalId="addSecret">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first secret</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
    </>
  )
}

export default SecretsPage
