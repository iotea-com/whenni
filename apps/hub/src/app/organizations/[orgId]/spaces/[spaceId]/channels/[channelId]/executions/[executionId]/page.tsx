import gruentClient from '@gruent/hub/lib/gruent'
import NodeExecutionLogEntry from '@gruent/libs/frontend/components/molecules/NodeExecutionLogEntry'
import Link from 'next/link'
import dayjs from 'dayjs'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Channel execution details | GRUENT',
}

const ChannelExecutionByIdPage = async ({ params }) => {
  const { orgId, spaceId, channelId, executionId: channelExecutionId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: channelExecution, errors: getChannelExecutionErrors } = await gruentClient(
    accessToken,
  ).channels.executions.get(spaceId, channelExecutionId)

  if (getChannelExecutionErrors && getChannelExecutionErrors.length > 0)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Could not retrieve the channel execution"
          description={`${getChannelExecutionErrors[0]}`}
          variant="error"
        />
      </div>
    )

  if (!channelExecution)
    return (
      <div className="mt-4 px-8 py-2 max-w-full mb-10">
        <Callout
          className="max-w-xl"
          title="Invalid channel execution"
          description={`This channel execution does not exist. Try logging out and logging back in.`}
          variant="warning"
        />
      </div>
    )

  const { executionId, status, startTime, endTime, durationMs } = channelExecution

  const nodeExecutionLogsJsx = (() => {
    const { nodeExecutionLogs } = channelExecution

    return nodeExecutionLogs.map((nodeExecutionLog, index) => {
      const { level, timestamp, logAttributes } = nodeExecutionLog
      const { node } = logAttributes

      return (
        <NodeExecutionLogEntry
          key={`${node}-${level}-${timestamp}-${index}`}
          nodeExecutionLog={nodeExecutionLog}
        />
      )
    })
  })()

  return (
    <div className="mt-4 px-8 py-2 max-w-full mb-10">
      <div className="mb-4 text-gray-600 dark:text-gray-500">
        <small>
          <Link href={`/organizations/${orgId}/spaces/${spaceId}`}>Space</Link> &gt;&nbsp;
          <Link href={`/organizations/${orgId}/spaces/${spaceId}/channels/${channelId}`}>
            Channel
          </Link>{' '}
          &gt;&nbsp;
          {channelExecutionId}
        </small>
      </div>
      <h1>{executionId}</h1>
      {status === 'COMPLETED' || status === 'ERRORED' || status === 'TIMEOUT' ? (
        <p>
          <small>
            {status} - Started at {dayjs(startTime).format('MM/DD/YYYY HH:mm:ss.SSS')}, ended at{' '}
            {dayjs(endTime).format('MM/DD/YYYY HH:mm:ss.SSS')} ({durationMs}ms)
          </small>
        </p>
      ) : (
        <p>
          <small>{status}</small>
        </p>
      )}
      <div style={{ width: '100%', maxWidth: 'calc(100vw - 150px)' }}>{nodeExecutionLogsJsx}</div>
    </div>
  )
}

export default ChannelExecutionByIdPage
