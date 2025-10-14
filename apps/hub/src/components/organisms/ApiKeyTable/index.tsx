'use client'

import { FC, useState } from 'react'
import {
  ColumnDef,
  createColumnHelper,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { ApiKey, PermissionSet } from '@prisma/client'
import ListTable from '@iotea/hub/components/atoms/ListTable'
import ApiKeyContextMenu from './ApiKeyContextMenu'
import {
  RemixIcon,
  riClipboardFill,
  riInformationLine,
  riKeyLine,
  riListSettingsLine,
} from '@mwarnerdotme/react-remixicon'
import copyToClipboard from '@iotea/libs/frontend/util/copyToClipboard'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'
// import IndeterminateCheckbox from "../../atoms/IndeterminateCheckbox"

type Props = {
  apiKeys: ApiKey[]
  orgId: string
  spaceId?: string
  permissionSets: PermissionSet[]
  page: number
  totalPages: number
  totalResults: number
}

const ApiKeysTable: FC<Props> = ({
  apiKeys,
  orgId,
  spaceId,
  permissionSets,
  page,
  totalPages,
  totalResults,
}) => {
  const [rowSelection, setRowSelection] = useState({})

  const router = useRouter()

  const columnHelper = createColumnHelper<ApiKey>()

  const columns: ColumnDef<ApiKey>[] = [
    // TODO: enable selection when there is bulk action support in the UI for selected rows
    // {
    //   id: 'select',
    //   cell: ({ row }) => (
    //     <IndeterminateCheckbox
    //       {...{
    //         checked: row.getIsSelected(),
    //         disabled: !row.getCanSelect(),
    //         indeterminate: row.getIsSomeSelected(),
    //         onChange: row.getToggleSelectedHandler(),
    //       }}
    //     />
    //   ),
    // },
    columnHelper.accessor('name', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riInformationLine} />
          <span>Name</span>
        </div>
      ),
      cell: (info) => info.getValue(),
      enableSorting: true,
    }),
    columnHelper.accessor('id', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riKeyLine} />
          <span>Key</span>
        </div>
      ),
      cell: (info) => {
        const apiKey = info.getValue()
        return (
          <div className="flex gap-2">
            <span>&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;&#8226;</span>
            <RemixIcon
              icon={riClipboardFill}
              className="cursor-pointer"
              onClick={async () => {
                const { error } = await copyToClipboard(apiKey)
                if (error)
                  addToast({
                    title: 'Could not copy text',
                    body: 'Could not copy the API key to your clipboard. This may be due to a browser permission issue.',
                  })
                else
                  addToast({
                    title: 'Copied text',
                    body: 'The API key should now be in your clipboard!',
                  })
              }}
            />
          </div>
        )
      },
      enableSorting: true,
      size: 300,
    }),
    columnHelper.accessor('organizationPermissionSetId', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riListSettingsLine} />
          <span>Permission Set</span>
        </div>
      ),
      cell: (info) => {
        const permissionSet = permissionSets.find((ps) => {
          return ps.id === info.getValue()
        })

        if (!permissionSet) return info.getValue()
        return permissionSet.name
      },
      enableSorting: true,
      size: 300,
    }),
    columnHelper.display({
      id: 'contextMenu',
      cell: (info) => (
        <ApiKeyContextMenu orgId={orgId} spaceId={spaceId} apiKeyId={info.row.original.id} />
      ),
      size: 25,
      meta: {
        align: 'right',
      },
    }),
  ]

  const table = useReactTable<ApiKey>({
    data: apiKeys,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getRowId: (original) => original.id,
    enableRowSelection: true,
    manualPagination: true,
    state: {
      rowSelection,
    },
    onRowSelectionChange: setRowSelection,
  })

  const handlePageChange = (newPage: number) => {
    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.set('apiKeysTablePage', newPage.toString())
    router.replace(newLocation.toString(), { scroll: false })
  }

  return (
    <ListTable
      title="API Keys"
      description={`Programmatically access the IOTEA platform.`}
      table={table}
      page={page}
      totalPages={totalPages}
      handlePageChange={handlePageChange}
    >
      {totalResults === 1 ? (
        <small>
          {totalResults} API key in this {spaceId ? 'space' : 'organization'}
        </small>
      ) : (
        <small>
          {totalResults} API keys in this {spaceId ? 'space' : 'organization'}
        </small>
      )}
    </ListTable>
  )
}

export default ApiKeysTable
