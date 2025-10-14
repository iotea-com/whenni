import notion from '@iotea/site/lib/notion'
import NotionBlockRenderer from '@iotea/libs/frontend/components/templates/NotionBlockRenderer'
import type { Metadata, ResolvingMetadata } from 'next'
import Link from 'next/link'
import {
  BlockObjectResponse,
  PageObjectResponse,
  UserObjectResponse,
} from '@notionhq/client/build/src/api-endpoints'
import dayjs from 'dayjs'
import {
  calculateBlogPostSlug,
  slugToNotionQuery,
} from '@iotea/libs/frontend/util/calculateBlogPostSlug'
import styles from './page.module.scss'
import { notFound } from 'next/navigation'
import TableOfContents from './TableOfContents'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { RemixIcon, riGitForkLine, riPriceTagLine } from '@mwarnerdotme/react-remixicon'

import Image from 'next/image'
import redis from '@iotea/site/lib/upstash'

type MetadataCache = {
  pageId: string
  title?: string
  description?: string
  coverImageUrl?: string
  publishedTime?: string
  authorImageUrl?: string
  authorName?: string
}

const CacheTime = 60 * 60 * 24 // 1 day

const buildMetadataCache = async (slug: string) => {
  // Check if the metadata cache already exists - return it if it does
  const exists = await redis.exists(`metadata-${slug}`)
  if (exists) {
    const metadataCache = (await redis.get(`metadata-${slug}`)) as MetadataCache
    return metadataCache
  }

  // First get pages that might match by doing a contains filter on the title
  const slugQuery = slugToNotionQuery(slug)

  const response = await notion.dataSources.query({
    data_source_id: process.env.NOTION_BLOG_DATABASE_ID!,
    filter: {
      property: 'Title',
      title: {
        contains: slugQuery,
      },
    },
  })

  // Get pages that might match by doing a contains filter on the title
  const page = response.results.find((page: PageObjectResponse) => {
    const title = page.properties.Title
    if (title.type === 'title') {
      const pageSlug = calculateBlogPostSlug(title.title[0].plain_text)
      return pageSlug === slug
    }
    return false
  }) as PageObjectResponse

  if (!page) {
    throw new Error(`Page with slug ${slug} not found`)
  }

  // Calculate the metadata cache
  const title = (() => {
    if (page.properties.Title.type === 'title') return page.properties.Title.title[0].plain_text
    return undefined
  })()

  const description = (() => {
    if (page.properties.Overview.type === 'rich_text')
      return page.properties.Overview.rich_text[0].plain_text
    return undefined
  })()

  const coverImageUrl = (() => {
    if (page.cover && page.cover.type === 'external') return page.cover.external.url
    if (page.cover && page.cover.type === 'file') return page.cover.file.url
    return undefined
  })()

  const authorImageUrl = (() => {
    return page.properties.Author.type === 'people' &&
      (page.properties.Author.people[0] as UserObjectResponse).type === 'person' &&
      (page.properties.Author.people[0] as UserObjectResponse).avatar_url
      ? ((page.properties.Author.people[0] as UserObjectResponse).avatar_url ??
          '/img/logos/app-icon-primary.png')
      : '/img/logos/app-icon-primary.png'
  })()

  const authorName = (() => {
    return page.properties.Author.type === 'people' &&
      (page.properties.Author.people[0] as UserObjectResponse).name
      ? ((page.properties.Author.people[0] as UserObjectResponse).name ?? 'IOTEA')
      : 'IOTEA'
  })()

  const metadataCache: MetadataCache = {
    pageId: page.id,
    title,
    description,
    coverImageUrl,
    publishedTime: dayjs(
      page.properties['Publish Date'].type === 'date'
        ? page.properties['Publish Date'].date?.start
        : null,
    ).toISOString(),
    authorImageUrl,
    authorName,
  }

  // Cache the metadata for 1 hour
  await redis.set(`metadata-${slug}`, metadataCache, {
    ex: CacheTime,
  })

  return metadataCache
}

export async function generateMetadata({ params }, _parent: ResolvingMetadata): Promise<Metadata> {
  const { slug } = await params

  try {
    const metadataCache = await buildMetadataCache(slug)
    return {
      title: metadataCache.title,
      description: metadataCache.description,
      openGraph: {
        type: 'article',
        url: `https://iotea.com/blog/${slug}`,
        title: metadataCache.title,
        description: metadataCache.description,
        images: metadataCache.coverImageUrl ? [metadataCache.coverImageUrl] : [],
        publishedTime: metadataCache.publishedTime,
      },
    }
  } catch (err) {
    throw new Error(err)
  }
}

const BlogPostPage = async ({ params }) => {
  const { slug } = await params

  const metadataCache = await (async () => {
    const cache = await redis.get(`metadata-${slug}`)
    if (cache) return cache as MetadataCache

    try {
      const newCache = await buildMetadataCache(slug)
      return newCache
    } catch (_err) {
      return null
    }
  })()

  if (!metadataCache) notFound()

  const fetchBlocksPage = async (blockId: string, cursor?: string) => {
    const { results, next_cursor } = await notion.blocks.children.list({
      block_id: blockId,
      start_cursor: cursor,
      page_size: 100,
    })
    return { results: results as BlockObjectResponse[], next_cursor }
  }

  const fetchBlockWithChildren = async (block: BlockObjectResponse) => {
    if (!block.has_children) return block

    const children = await fetchAllBlocks(block.id)
    return { ...block, children }
  }

  const fetchAllBlocks = async (blockId: string): Promise<BlockObjectResponse[]> => {
    // Get cached blocks if they exist
    if (await redis.exists(`blocks-${slug}`))
      return (await redis.get(`blocks-${slug}`)) as BlockObjectResponse[]

    // Fetch all blocks and children if they don't exist in cache
    const accumulator = async (
      acc: BlockObjectResponse[],
      cursor?: string,
    ): Promise<BlockObjectResponse[]> => {
      const { results, next_cursor } = await fetchBlocksPage(blockId, cursor)

      const blocksWithChildren = await Promise.all(results.map(fetchBlockWithChildren))

      const newAcc = [...acc, ...blocksWithChildren]
      return next_cursor ? accumulator(newAcc, next_cursor) : newAcc
    }

    return accumulator([])
  }

  const blocks = await fetchAllBlocks(metadataCache.pageId)
  await redis.set(`blocks-${slug}`, blocks, {
    ex: CacheTime,
  })

  // Create nested table of contents
  const tableOfContents = blocks
    .filter((block) => block.type === 'heading_2' || block.type === 'heading_3')
    .reduce(
      (toc, block) => {
        const headingText = block[block.type].rich_text[0].text.content

        if (block.type === 'heading_2') {
          toc[headingText] = {
            type: 'heading_2',
            subheadings: {},
          }
        } else if (block.type === 'heading_3') {
          // Find the last h2 heading
          const lastH2 = Object.keys(toc).pop()
          if (lastH2) {
            toc[lastH2].subheadings[headingText] = {
              type: 'heading_3',
            }
          }
        }
        return toc
      },
      {} as {
        [key: string]: {
          type: 'heading_2'
          subheadings: {
            [key: string]: {
              type: 'heading_3'
            }
          }
        }
      },
    )

  return (
    <>
      <div className="grid grid-cols-12 lg:max-w-4xl mx-auto min-h-screen">
        <TableOfContents
          tableOfContents={tableOfContents}
          className="hidden md:block md:sticky top-5 mt-10"
        />
        <div className="col-span-12 md:col-span-8 lg:col-span-9">
          <nav className="mx-auto mb-4 text-sm">
            <Link className="text-gray-600 transition underline hover:text-green-700" href="/blog">
              Blog
            </Link>{' '}
            / {metadataCache.title}
          </nav>
          <article id={styles.article} className="mx-auto">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={metadataCache.coverImageUrl ?? ''}
              alt="article thumbnail"
              className="w-full mb-4 shadow-lg rounded"
            />
            <h1>{metadataCache.title}</h1>
            <div className="flex items-center">
              <div className="rounded-full overflow-hidden">
                <Image
                  src={metadataCache.authorImageUrl ?? ''}
                  alt="author's thumbnail"
                  height={45}
                  width={45}
                />
              </div>
              <div className="ml-2 flex flex-col">
                <p className="text-gray-700">{metadataCache.authorName}</p>
                <p>
                  <span className="text-gray-600">
                    {dayjs(metadataCache.publishedTime).format('MMM DD, YYYY')}
                  </span>
                </p>
              </div>
            </div>
            <TableOfContents tableOfContents={tableOfContents} className="md:hidden mt-4" />
            <div className={styles.content}>
              <NotionBlockRenderer blocks={blocks} />
            </div>
          </article>
        </div>
      </div>
      <hr className="mt-8 mb-16 max-w-4xl mx-auto" />
      <div className="grid grid-cols-1 md:grid-cols-2 gap-10 max-w-4xl mx-auto">
        <div>
          <h2 className="text-2xl">Connect in minutes</h2>
          <p className="mt-4 mb-6 text-gray-500">
            Create an account and start connecting for free until you&apos;re ready to commit.
          </p>
          <div className="flex">
            <Button text="Start Now" modalId="betaSignup" />
            <a href="/pricing">
              <Button className="ml-4" variant="transparent" text="See Pricing" />
            </a>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
          <div>
            <RemixIcon icon={riGitForkLine} size="2x" className="text-green-500 rotate-180" />
            <h3 className="font-bold text-gray-700 mt-3 mb-1">Connect anything</h3>
            <p className="text-gray-500">
              Connect devices and services using any protocol with our platform.
            </p>
          </div>
          <div>
            <RemixIcon icon={riPriceTagLine} size="2x" className="text-green-500" />
            <h3 className="font-bold text-gray-700 mt-3 mb-1">Pay for usage</h3>
            <p className="text-gray-500">
              Stop spending more for team-based subscriptions. Pay only for what you use.
            </p>
          </div>
        </div>
      </div>
    </>
  )
}

export default BlogPostPage
