'use client'

import { FC, useState } from 'react'
import {
  ColumnDef,
  createColumnHelper,
  getCoreRowModel,
  getPaginationRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { OrganizationMember, PermissionSet, User } from '@prisma/client'
import ListTable from '@iotea/hub/components/atoms/ListTable'
import { useRouter } from 'next/navigation'
import MemberContextMenu from './MembersContextMenu'
import { riAdminLine, riMailLine, riUserSettingsLine } from '@mwarnerdotme/react-remixicon'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
// import IndeterminateCheckbox from "../../atoms/IndeterminateCheckbox"

type Member = OrganizationMember & {
  user: User
}

type Props = {
  users: Member[]
  permissionSets: PermissionSet[]
  orgId: string
  page: number
  totalPages: number
  totalResults: number
}

const MembersTable: FC<Props> = ({
  orgId,
  users,
  permissionSets,
  page,
  totalPages,
  totalResults,
}) => {
  const [rowSelection, setRowSelection] = useState({})

  const router = useRouter()

  const columnHelper = createColumnHelper<Member>()

  const columns: ColumnDef<Member>[] = [
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
    columnHelper.accessor('user.email', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riMailLine} />
          <span>Email</span>
        </div>
      ),
      cell: (info) => info.getValue(),
      enableSorting: true,
    }),
    columnHelper.accessor('role', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riAdminLine} />
          <span>Role</span>
        </div>
      ),
      cell: (info) => info.getValue(),
      enableSorting: true,
      size: 150,
    }),
    columnHelper.accessor('organizationPermissionSetId', {
      header: () => (
        <div className="flex items-center">
          <RemixIcon className="mr-1" icon={riUserSettingsLine} />
          <span>Permission Set</span>
        </div>
      ),
      cell: (info) => {
        if (info.row.original.role === 'ADMIN') return 'Admin'
        const permissionSet = permissionSets.find((ps) => {
          return ps.id === info.getValue()
        })

        if (!permissionSet) return info.getValue()

        return permissionSet.name
      },
      enableSorting: true,
      size: 200,
    }),
    columnHelper.display({
      id: 'contextMenu',
      cell: (info) => (
        <MemberContextMenu
          orgId={orgId}
          userId={info.row.original.userId}
          admin={info.row.original.role === 'ADMIN'}
          initialPermissionSetId={info.row.original.organizationPermissionSetId}
          permissionSets={permissionSets}
        />
      ),
      size: 25,
      meta: {
        align: 'right',
      },
    }),
  ]

  const table = useReactTable<Member>({
    data: users,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getRowId: (original) => original.userId,
    enableRowSelection: true,
    manualPagination: true,
    state: {
      rowSelection,
    },
    onRowSelectionChange: setRowSelection,
  })

  const handlePageChange = (newPage: number) => {
    const newLocation = new URL(window.location.toString())
    newLocation.searchParams.set('membersTablePage', newPage.toString())
    router.replace(newLocation.toString(), { scroll: false })
  }

  return (
    <ListTable
      title="Members"
      description="Invite your team to collaborate on projects."
      table={table}
      page={page}
      totalPages={totalPages}
      handlePageChange={handlePageChange}
    >
      {totalResults === 1 ? (
        <small>{totalResults} member in this organization</small>
      ) : (
        <small>{totalResults} members in this organization</small>
      )}
    </ListTable>
  )
}

export default MembersTable
