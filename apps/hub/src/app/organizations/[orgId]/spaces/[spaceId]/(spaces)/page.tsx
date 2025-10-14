import CreateChannelModal from '@iotea/hub/components/modals/CreateChannelModal'
import ioteaClient from '@iotea/hub/lib/iotea'
import ChannelsTable from '@iotea/hub/components/organisms/ChannelsTable'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import Container from '@iotea/libs/frontend/components/templates/Container'
import { riAddCircleLine, riGitForkLine } from '@mwarnerdotme/react-remixicon'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
import ListTablePlaceholder from '@iotea/hub/components/molecules/ListTablePlaceholder'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
import { Tag } from '@prisma/client'
const SpaceDashboardPage = async ({ params, searchParams }) => {
  const { orgId, spaceId } = await params
  const { page: requestedPage = 1, q: filter, tagFilter: tagFilterString } = await searchParams

  const tagFilter = tagFilterString ? tagFilterString.split(',') : undefined

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const {
    data: channels,
    errors,
    page,
    totalPages,
    totalResults,
  } = await ioteaClient(accessToken).channels.list(spaceId, {
    page: requestedPage,
    filter,
    tagFilter,
  })

  if (errors && errors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the channels for this space"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!channels) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Channel does not exist"
          description={`The channel does not exist. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  const { data: tagListResponse, errors: tagListErrors } =
    await ioteaClient(accessToken).tags.list(spaceId)

  if (tagListErrors && tagListErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the tags for this space"
          description={`${tagListErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!tagListResponse) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the tags for this space"
          description={`The tags could not be listed. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  const tags: Tag[] = []
  for (const [_, v] of Object.entries(tagListResponse)) {
    tags.push(v.tag)
  }

  return (
    <>
      <CreateChannelModal spaceId={spaceId} />
      <Container>
        <h1 className="text-2xl mb-6">Welcome to your space</h1>
      </Container>
      <Container>
        <section
          className="relative bg-gray-50 rounded py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          {((Number(totalResults) > 0 && !filter) || filter || tagFilter) && (
            <ChannelsTable
              channels={channels}
              tags={tags}
              orgId={orgId}
              spaceId={spaceId}
              page={page!}
              totalPages={totalPages!}
              totalResults={totalResults ?? 0}
              searchFilter={filter}
              tagFilter={tagFilter}
            />
          )}
          {Number(totalResults) <= 0 && !filter && !tagFilter && (
            <ListTablePlaceholder
              title="Channels"
              description="Build connections to power your ideas."
            >
              <RemixIcon icon={riGitForkLine} className="text-gray-600 mb-2 rotate-180" size="3x" />
              <p className="text-gray-600 text-center mb-2">No channels... yet!</p>
              <Button variant="primary" className="text-sm" modalId="createChannel">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first channel</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
    </>
  )
}

export default SpaceDashboardPage
