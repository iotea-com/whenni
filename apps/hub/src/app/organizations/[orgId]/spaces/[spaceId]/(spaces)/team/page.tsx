import ioteaClient from '@iotea/hub/lib/iotea'
import CreatePermissionSetModal from '@iotea/hub/components/modals/CreatePermissionSetModal'
import MembersTable from '@iotea/hub/components/organisms/MembersTable'
import ApiKeysTable from '@iotea/hub/components/organisms/ApiKeyTable'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import PermissionSetTable from '@iotea/hub/components/organisms/PermissionSetTable'
import SearchBar from '@iotea/hub/components/molecules/SearchBar'
import Container from '@iotea/libs/frontend/components/templates/Container'
import CreateApiKeyModal from '@iotea/hub/components/modals/CreateApiKeyModal'
import { riAddCircleLine, riTerminalBoxLine } from '@mwarnerdotme/react-remixicon'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import ListTablePlaceholder from '@iotea/hub/components/molecules/ListTablePlaceholder'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'

export const metadata = {
  title: 'Team | IOTEA',
}

const TeamPage = async ({ params, searchParams }) => {
  const { orgId, spaceId } = await params
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
    errors: listUsersErrors,
    page: membersPage,
    totalPages: membersTotalPages,
    totalResults: totalMembers,
  } = await ioteaClient(accessToken).organizations.members.list(orgId, {
    page: requestedMembersTablePage,
    filter: membersTableFilter,
  })

  if (listUsersErrors && listUsersErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the members for this space"
          description={`${listUsersErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!users) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the channels for this space"
          description={`The channels could not be retrieved. Try logging out and logging back in.`}
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
  } = await ioteaClient(accessToken).permissions.list(spaceId, {
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

  if (!permissionSets) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the permission sets for this space"
          description={`The permission sets could not be retrieved. Try logging out and logging back in.`}
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
  } = await ioteaClient(accessToken).apiKeys.list(
    orgId,
    {
      page: requestedApiKeysTablePage,
    },
    spaceId,
  )

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

  if (!apiKeys) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not get the API keys for this space"
          description={`The API keys could not be retrieved. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <>
      <CreateApiKeyModal orgId={orgId} spaceId={spaceId} permissionSets={permissionSets} />
      <CreatePermissionSetModal orgId={orgId} spaceId={spaceId} />
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          <div className="absolute top-8 right-10 mb-1 flex items-end gap-2">
            <div className="grow" />
            <SearchBar
              paramKey="mq"
              initialFilter={membersTableFilter}
              placeholder="Search members"
            />
          </div>
          <MembersTable
            orgId={orgId}
            users={users}
            permissionSets={permissionSets}
            page={membersPage!}
            totalPages={membersTotalPages!}
            totalResults={totalMembers ?? 0}
          />
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
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
                spaceId={spaceId}
                apiKeys={apiKeys}
                permissionSets={permissionSets}
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
              <RemixIcon icon={riTerminalBoxLine} className="text-gray-600 mb-2" size="3x" />
              <p className="text-gray-600 text-center mb-2">No API keys... yet!</p>
              <Button variant="primary" className="text-sm" modalId="addApiKey">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first API key</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
      <Container className="mb-4">
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
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
            spaceId={spaceId}
            permissionSets={permissionSets}
            page={permissionSetsPage!}
            totalPages={permissionSetsTotalPages!}
            totalResults={totalPermissionSets ?? 0}
          />
        </section>
      </Container>
    </>
  )
}

export default TeamPage
