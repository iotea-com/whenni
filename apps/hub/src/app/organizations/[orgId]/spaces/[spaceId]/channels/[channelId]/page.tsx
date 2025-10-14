import ioteaClient from '@iotea/hub/lib/iotea'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import ChannelExecutionsTable from '@iotea/hub/components/organisms/ChannelExecutionsTable'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Channel details | IOTEA',
}

const ChannelDetailsPage = async ({ params, searchParams }) => {
  const { orgId, spaceId, channelId } = await params
  const { page: requestedPage = 1 } = await searchParams

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: channel, errors: getChannelErrors } = await ioteaClient(accessToken).channels.get(
    spaceId,
    channelId,
  )

  if (getChannelErrors && getChannelErrors.length > 0)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not retrieve the channel"
          description={`${getChannelErrors[0]}`}
          variant="error"
        />
      </div>
    )

  const {
    data: channelExecutions,
    errors: getChannelExecutionsErrors,
    page,
    totalPages,
    totalResults,
  } = await ioteaClient(accessToken).channels.executions.list(spaceId, channelId, {
    page: requestedPage,
  })

  if (getChannelExecutionsErrors && getChannelExecutionsErrors.length > 0)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not retrieve the channel executions"
          description={`${getChannelExecutionsErrors[0]}`}
          variant="error"
        />
      </div>
    )

  const { data: channelStatus, errors: _channelStatusErrors } = await ioteaClient(
    accessToken,
  ).channels.status(spaceId, channelId)

  if (!channel)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Invalid channel"
          description={`This channel does not exist. Try logging out and logging back in.`}
          variant="warning"
        />
      </div>
    )

  return (
    <div className="mt-4 px-8 py-2 max-w-full mb-10">
      <h1>{channel.name}</h1>
      <p>
        <small>Monitor and edit channels</small>
      </p>
      <Button
        className="mt-2"
        href={`/organizations/${orgId}/spaces/${spaceId}/channels/${channelId}/edit`}
        text="Edit"
      />
      <h2 className="mt-8">Channel Execution Logs</h2>
      {channelStatus?.status === 'PUBLISHED' && channel.publishedAt && (
        <p>
          <small>
            Channel was published on {new Date(channel.publishedAt).toLocaleDateString()} at{' '}
            {new Date(channel.publishedAt).toLocaleTimeString()}.
          </small>
        </p>
      )}
      {channelStatus?.status === 'PUBLISHED' && !channel.publishedAt && (
        <p>
          <small>Channel is published.</small>
        </p>
      )}
      {channelStatus?.status === 'INITIALIZING' && (
        <p>
          <small>Channel is initializing.</small>
        </p>
      )}
      {channelStatus?.status === 'UNPUBLISHED' && (
        <p>
          <small>Channel is not published.</small>
        </p>
      )}
      <div className="max-w-4xl">
        <ChannelExecutionsTable
          orgId={orgId}
          spaceId={spaceId}
          channelExecutionList={channelExecutions ?? []}
          page={page ?? 1}
          totalPages={totalPages ?? 1}
          totalResults={totalResults ?? 0}
        />
      </div>
    </div>
  )
}

export default ChannelDetailsPage
