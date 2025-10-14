'use client'

import { FC, useState } from 'react'
import dayjs from 'dayjs'
import { RemixIcon, riArrowDownSLine, riArrowUpSLine } from '@mwarnerdotme/react-remixicon'

type LogAttributes = {
  execution_id: string
  node: string
  level: string
  time: string
  data?: string
  'base64-data'?: string
}

export type NodeExecutionLog = {
  traceId: string
  spanId: string
  timestamp: string
  level: string
  body: string
  logAttributes: LogAttributes
  resourceAttributes: Record<string, string>
}

const NodeExecutionLogEntry: FC<{
  nodeExecutionLog: NodeExecutionLog
}> = ({ nodeExecutionLog }) => {
  const [open, setOpen] = useState(false)

  const { timestamp, body, level, logAttributes } = nodeExecutionLog
  const { node } = logAttributes

  const borderColor = (() => {
    switch (level) {
      case 'error':
        return 'border-error'
      case 'warn':
        return 'border-warning'
      case 'info':
      default:
        return 'border-gray-300 hover:border-blue-500'
    }
  })()

  const stringifiedMessage = (() => {
    if (logAttributes.data) {
      const logData = (() => {
        try {
          return JSON.parse(logAttributes.data)
        } catch (_e) {
          // eslint-disable-line
          return logAttributes.data
        }
      })()
      return `${nodeExecutionLog.body}\n\n${JSON.stringify(logData, null, 2)}`
    }

    if (logAttributes['base64-data']) {
      const logData = (() => {
        try {
          const d = decodeURIComponent(atob(logAttributes['base64-data']))
          return JSON.parse(d)
        } catch (_e) {
          // eslint-disable-line
          return logAttributes['base64-data']
        }
      })()
      return `${nodeExecutionLog.body}\n\n${JSON.stringify(logData, null, 2)}`
    }

    return body
  })()

  return (
    <div className={`mt-2 rounded transition border ${borderColor} w-full`}>
      <div
        key={nodeExecutionLog.timestamp}
        className={`transition ${open ? 'border-b' : ''} ${borderColor} p-4 cursor-pointer flex`}
        onClick={() => setOpen((current) => !current)}
      >
        <p className="select-none">
          {dayjs(timestamp).format('MM-DD-YYYY HH:mm:ss.SSS ')} - {node} ({level})
        </p>
        <div className="grow" />
        {open ? (
          <RemixIcon className="text-gray-500" icon={riArrowUpSLine} />
        ) : (
          <RemixIcon className="text-gray-500" icon={riArrowDownSLine} />
        )}
      </div>
      {open && (
        <section className="max-w-full overflow-x-scroll">
          <pre className="whitespace-pre-wrap dark:text-gray-200 pt-4 px-4 pb-1 text-sm">
            {stringifiedMessage}
          </pre>
        </section>
      )}
    </div>
  )
}

export default NodeExecutionLogEntry
