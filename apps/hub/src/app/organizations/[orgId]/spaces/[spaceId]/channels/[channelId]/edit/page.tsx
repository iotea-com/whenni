import getAccessToken from '@gruent/hub/util/getAccessToken'
import ChannelByIdEditClientPage from './client'
import gruentClient from '@gruent/hub/lib/gruent'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Edit channel | GRUENT',
}

const ChannelByIdEditPage = async ({ params }) => {
  const { orgId, spaceId, channelId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: space, errors: spaceErrors } = await gruentClient(accessToken).spaces.get(
    orgId,
    spaceId,
  )

  if (spaceErrors && spaceErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the space details"
          description={`${spaceErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  const { data: organization, errors: organizationErrors } = await gruentClient(
    accessToken,
  ).organizations.get(space?.organizationId ?? '')

  if (organizationErrors && organizationErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the organization details"
          description={`${organizationErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <ChannelByIdEditClientPage
      orgId={orgId}
      spaceId={spaceId}
      channelId={channelId}
      organization={organization}
      space={space}
    />
  )
}

export default ChannelByIdEditPage
