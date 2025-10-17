'use client'

import { FC, useRef, useState } from 'react'
import {
  ColumnDef,
  createColumnHelper,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from '@tanstack/react-table'
import ListTable from '@iotea/hub/components/atoms/ListTable'
import Link from 'next/link'
import dayjs from 'dayjs'
import { useRouter } from 'next/navigation'
import ListTablePlaceholder from '../../molecules/ListTablePlaceholder'
import { RemixIcon, riEqualizerLine, riSettings6Line } from '@mwarnerdotme/react-remixicon'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { useClickOutside } from '@react-hooks-library/core'
// import IndeterminateCheckbox from "../../atoms/IndeterminateCheckbox"

type ChannelExecutionListItem = {
  executionId: string
  status: string
  startTime: string
  endTime: string
  durationMs: number
  channelId: string
  logAttributes: Record<string, string>
}

type Props = {
  channelExecutionList: ChannelExecutionListItem[]
  orgId: string
  spaceId: string
  page: number
  totalPages: number
  totalResults: number
  autoRefresh?: number
  isLoading?: boolean
  onPageChange?: (newPage: number) => void
  onExecutionClick?: (executionId: string) => void
  statusFilter?: string
  onStatusFilterChange?: (status: string | null) => void
}

const ChannelExecutionsTable: FC<Props> = ({
  orgId,
  spaceId,
  channelExecutionList,
  page,
  totalPages,
  totalResults,
  autoRefresh,
  isLoading,
  onPageChange,
  onExecutionClick,
  statusFilter,
  onStatusFilterChange,
}) => {
  const [rowSelection, setRowSelection] = useState({})
  const [showFilterOptions, setShowFilterOptions] = useState(false)
  const router = useRouter()
  const filterOptionsRef = useRef<HTMLDivElement>(null)

  useClickOutside(filterOptionsRef, () => setShowFilterOptions(false))

  const columnHelper = createColumnHelper<ChannelExecutionListItem>()

  const columns: ColumnDef<ChannelExecutionListItem>[] = [
    columnHelper.accessor('executionId', {
      header: () => 'Execution ID',
      cell: (info) => {
        if (onExecutionClick) {
          return (
            <span
              className="cursor-pointer"
              onClick={() => {
                onExecutionClick?.(info.getValue())
              }}
            >
              {info.getValue()}
            </span>
          )
        }

        return (
          <Link
            href={`/organizations/${orgId}/spaces/${spaceId}/channels/${info.row.original.channelId}/executions/${info.row.id}`}
          >
            {info.getValue()}
          </Link>
        )
      },
      enableSorting: true,
    }),
    columnHelper.accessor('status', {
      header: () => 'Status',
      cell: (info) => info.getValue(),
      enableSorting: true,
    }),
    columnHelper.accessor('startTime', {
      header: () => 'Started At',
      cell: (info) => dayjs(info.getValue()).format('MM/DD/YYYY HH:mm:ss.SSS'),
      enableSorting: true,
    }),
    columnHelper.display({
      id: 'duration',
      header: () => 'Duration',
      cell: (info) => {
        const { durationMs } = info.row.original

        const durationString = (() => {
          if (durationMs < 500) return `${durationMs}ms`
          return `${durationMs / 1000}s`
        })()

        return durationString
      },
      enableSorting: true,
    }),
  ]

  const table = useReactTable<ChannelExecutionListItem>({
    data: channelExecutionList,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getRowId: (original) => original.executionId,
    enableRowSelection: true,
    manualPagination: true,
    state: {
      rowSelection,
    },
    onRowSelectionChange: setRowSelection,
  })

  const handlePageChange = (newPage: number) => {
    if (onPageChange) {
      onPageChange(newPage)
      return
    }

    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.set('page', newPage.toString())
    router.replace(newLocation.toString(), { scroll: false })
  }

  const handleStatusFilterChange = (status: string) => {
    if (onStatusFilterChange) {
      onStatusFilterChange(status)
      return
    }

    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.set('statusFilter', status)
    router.replace(newLocation.toString(), { scroll: false })
  }

  const clearFilters = () => {
    if (onStatusFilterChange) {
      onStatusFilterChange(null)
      return
    }

    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.delete('statusFilter')
    router.replace(newLocation.toString(), { scroll: false })
  }

  return (
    <>
      {((Number(totalResults) > 0 && !statusFilter) || statusFilter) && (
        <ListTable
          title="Channel Execution Logs"
          table={table}
          page={page}
          totalPages={totalPages}
          handlePageChange={handlePageChange}
          autoRefresh={autoRefresh}
          isLoading={isLoading}
          headerContent={
            <>
              <div className="flex flex-wrap items-end gap-2 mb-1 relative">
                <div className="grow" />
                <Button
                  variant="transparent"
                  className="px-2"
                  onClick={() => setShowFilterOptions((current) => !current)}
                >
                  <RemixIcon icon={riEqualizerLine} size="lg" />
                </Button>
                {showFilterOptions && (
                  <div
                    ref={filterOptionsRef}
                    className="absolute max-w-[500px] z-10 top-10 right-0 bg-gray-100 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-md p-4 flex flex-col gap-2 overflow-hidden"
                  >
                    <div className="border-b border-gray-200 dark:border-gray-800 pb-2">
                      <h2 className="text-sm font-medium mb-1">Tags</h2>
                      <div className="flex flex-wrap items-center gap-2">
                        {['info', 'error'].map((status) => {
                          if (statusFilter && statusFilter.includes(status))
                            return (
                              <Button
                                key={status}
                                variant="transparent"
                                className="text-sm border-primary!"
                                onClick={() => handleStatusFilterChange(status)}
                              >
                                {status}
                              </Button>
                            )

                          return (
                            <Button
                              key={status}
                              variant="transparent"
                              className="text-sm"
                              onClick={() => handleStatusFilterChange(status)}
                            >
                              {status}
                            </Button>
                          )
                        })}
                      </div>
                    </div>
                    <span
                      className="cursor-pointer transition text-xs text-gray-600 hover:text-red-500"
                      onClick={clearFilters}
                    >
                      Clear filters
                    </span>
                  </div>
                )}
              </div>
            </>
          }
        >
          {totalResults === 1 ? (
            <small>{totalResults} execution of this channel</small>
          ) : (
            <small>{totalResults} executions of this channel</small>
          )}
        </ListTable>
      )}
      {Number(totalResults) <= 0 && !statusFilter && (
        <ListTablePlaceholder
          title="Channel Executions"
          description="No executions found for this channel"
        >
          <RemixIcon icon={riSettings6Line} size="3x" />
          <p className="text-gray-500 px-4">No executions found for this channel... yet!</p>
        </ListTablePlaceholder>
      )}
    </>
  )
}

export default ChannelExecutionsTable
