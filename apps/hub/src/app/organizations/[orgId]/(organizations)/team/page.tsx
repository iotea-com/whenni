import CreateMemberModal from '@iotea/hub/components/modals/CreateMemberModal'
import ioteaClient from '@iotea/hub/lib/iotea'
import CreatePermissionSetModal from '@iotea/hub/components/modals/CreatePermissionSetModal'
import CreateApiKeyModal from '@iotea/hub/components/modals/CreateApiKeyModal'
import MembersTable from '@iotea/hub/components/organisms/MembersTable'
import PermissionSetTable from '@iotea/hub/components/organisms/PermissionSetTable'
import ApiKeysTable from '@iotea/hub/components/organisms/ApiKeyTable'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import SearchBar from '@iotea/hub/components/molecules/SearchBar'
import { riAddCircleLine, riTerminalBoxLine } from '@mwarnerdotme/react-remixicon'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Container from '@iotea/libs/frontend/components/templates/Container'
import ListTablePlaceholder from '@iotea/hub/components/molecules/ListTablePlaceholder'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Team | IOTEA',
}

const OrganizationTeamPage = async ({ params, searchParams }) => {
  const { orgId } = await params
  const {
    membersTablePage: requestedMembersTablePage = 1,
    apiKeysTablePage: requestedApiKeysTablePage = 1,
    permissionsTablePage: requestedPermissionsTablePage = 1,
    mq: membersTableFilter = '',
  } = await searchParams

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  // Get members
  const {
    data: users,
    errors: listMembersErrors,
    page: membersPage,
    totalPages: membersTotalPages,
    totalResults: totalMembers,
  } = await ioteaClient(accessToken).organizations.members.list(orgId, {
    page: requestedMembersTablePage,
    filter: membersTableFilter,
  })

  if (listMembersErrors && listMembersErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the members for this organization"
          description={`${listMembersErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  // Get permission sets
  const {
    data: permissionSets,
    errors: listPermissionSetsErrors,
    page: permissionSetsPage,
    totalPages: permissionSetsTotalPages,
    totalResults: totalPermissionSets,
  } = await ioteaClient(accessToken).permissions.list(orgId, {
    page: requestedPermissionsTablePage,
  })

  if (listPermissionSetsErrors && listPermissionSetsErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the permission sets for this space"
          description={`${listPermissionSetsErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  // Get API keys
  const {
    data: apiKeys,
    errors: listApiKeysErrors,
    page: apiKeysPage,
    totalPages: apiKeysTotalPages,
    totalResults: totalApiKeys,
  } = await ioteaClient(accessToken).apiKeys.list(orgId, {
    page: requestedApiKeysTablePage,
  })

  if (listApiKeysErrors && listApiKeysErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the API keys for this space"
          description={`${listApiKeysErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <CreateMemberModal orgId={orgId} />
      <CreatePermissionSetModal orgId={orgId} />
      <CreateApiKeyModal orgId={orgId} permissionSets={permissionSets ?? []} />
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px rgba(0,0,0,.03)' }}
        >
          <div className="absolute top-8 right-10 mb-1 flex items-end gap-2">
            <div className="grow" />
            <SearchBar
              paramKey="mq"
              initialFilter={membersTableFilter}
              placeholder="Search members"
            />
            <Button variant="transparent" className="text-sm" modalId="addMember">
              <RemixIcon className="mr-1" icon={riAddCircleLine} />
              <span>Add</span>
            </Button>
          </div>
          <MembersTable
            orgId={orgId}
            users={users ?? []}
            permissionSets={permissionSets ?? []}
            page={membersPage!}
            totalPages={membersTotalPages!}
            totalResults={totalMembers ?? 0}
          />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px rgba(0,0,0,.03)' }}
        >
          {Number(totalApiKeys) > 0 && (
            <>
              <div className="absolute top-8 right-10 mb-1 flex items-end gap-2">
                <div className="grow" />
                <Button variant="transparent" className="text-sm" modalId="addApiKey">
                  <RemixIcon className="mr-1" icon={riAddCircleLine} />
                  <span>Create</span>
                </Button>
              </div>
              <ApiKeysTable
                orgId={orgId}
                apiKeys={apiKeys ?? []}
                permissionSets={permissionSets ?? []}
                page={apiKeysPage!}
                totalPages={apiKeysTotalPages!}
                totalResults={totalApiKeys ?? 1}
              />
            </>
          )}
          {Number(totalApiKeys) <= 0 && (
            <ListTablePlaceholder
              title="API Keys"
              description="Programmatically access the IOTEA platform."
            >
              <RemixIcon
                icon={riTerminalBoxLine}
                className="text-gray-600 dark:text-gray-200 mb-2"
                size="3x"
              />
              <p className="text-gray-600 dark:text-gray-400 text-center mb-2">
                No API keys... yet!
              </p>
              <Button variant="primary" className="text-sm" modalId="addApiKey">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first API key</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
      <Container>
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px rgba(0,0,0,.03)' }}
        >
          <div className="absolute top-8 right-10 mb-1 flex items-end gap-2">
            <div className="grow" />
            <Button variant="transparent" className="text-sm" modalId="addPermissionSet">
              <RemixIcon className="mr-1" icon={riAddCircleLine} />
              <span>Create</span>
            </Button>
          </div>
          <PermissionSetTable
            orgId={orgId}
            permissionSets={permissionSets ?? []}
            page={permissionSetsPage!}
            totalPages={permissionSetsTotalPages!}
            totalResults={totalPermissionSets ?? 0}
          />
        </section>
      </Container>
    </>
  )
}

export default OrganizationTeamPage
