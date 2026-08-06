import notion from '@gruent/site/lib/notion'
import {
  PageObjectResponse,
  QueryDataSourceResponse,
  UserObjectResponse,
} from '@notionhq/client/build/src/api-endpoints'
import Image from 'next/image'
import Link from 'next/link'
import { calculateBlogPostSlug } from '@gruent/libs/frontend/util/calculateBlogPostSlug'

import styles from './page.module.scss'
import { Metadata } from 'next'
import dayjs from 'dayjs'
import redis from '@gruent/site/lib/upstash'

const CacheTime = 60 * 60 * 24 // 1 day

export const metadata: Metadata = {
  title: 'Updates, announcements, and guides | GRUENT Blog',
  description: 'Explore the latest developments in the internet of things in the GRUENT blog.',
  openGraph: {
    type: 'website',
    url: `https://gruent.com/blog`,
    title: 'Updates, announcements, and guides | GRUENT Blog',
    description: 'Explore the latest developments in the internet of things in the GRUENT blog.',
    images: ['https://gruent.com/img/logos/app-icon-primary.png'],
  },
}

const getTags = async () => {
  if (await redis.exists(`tags`)) return (await redis.get(`tags`)) as QueryDataSourceResponse

  const tagsResults = await notion.dataSources.query({
    data_source_id: process.env.NOTION_BLOG_TAG_DATABASE_ID!,
  })

  await redis.set(`tags`, tagsResults, {
    ex: CacheTime,
  })

  return tagsResults
}

const getPosts = async (
  tag: string,
  postCountByTag: Map<string, { tagId: string; tagName: string; postCount: number }>,
) => {
  if (await redis.exists(`posts-${tag}`))
    return (await redis.get(`posts-${tag}`)) as QueryDataSourceResponse

  const postsResults = await notion.dataSources.query({
    data_source_id: process.env.NOTION_BLOG_DATABASE_ID!,
    sorts: [
      {
        property: 'Publish Date',
        direction: 'descending',
      },
    ],
    filter: {
      and: [
        {
          property: 'Publish Date',
          date: {
            before: new Date().toISOString(),
          },
        },
        {
          property: 'Tags',
          relation: {
            contains: postCountByTag.get(tag)?.tagId || '',
          },
        },
      ],
    },
    page_size: 10,
  })

  await redis.set(`posts-${tag}`, postsResults, {
    ex: CacheTime,
  })

  return postsResults
}

const BlogPage = async ({ searchParams }) => {
  const { tag: tagParam } = await searchParams
  const tag = (tagParam as string) ?? 'All'

  // Retrieve tags from Notion
  const tagsResults = await getTags()

  // Process tags results to map post counts by tag
  const postCountByTag = new Map<string, { tagId: string; tagName: string; postCount: number }>(
    tagsResults.results.map((tag) => {
      const tagId = tag.id
      // @ts-expect-error custom property
      const tagName = tag.properties['Tag'].title[0].text.content
      // @ts-expect-error custom property
      const postCount = tag.properties['Live Count'].formula.number
      return [
        tagName,
        {
          tagId,
          tagName,
          postCount,
        },
      ]
    }),
  )

  // Retrieve posts from Notion
  const postsResults = await getPosts(tag, postCountByTag)

  const posts = postsResults.results as PageObjectResponse[]

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-4xl mb-3">GRUENT Blog</h1>
      <nav className="flex gap-4 mb-6">
        {Array.from(postCountByTag.values())
          .sort((a, b) => {
            // Ensure "All" is always at the front
            if (a.tagName === 'All') return -1
            if (b.tagName === 'All') return 1

            // Compare the rest alphabetically
            return a.tagName.localeCompare(b.tagName)
          })
          .map((t) => {
            const { tagName, postCount } = t
            if (postCount <= 0) return null

            if (tagName === tag)
              return (
                <Link
                  className="text-sm px-3 py-1 rounded-full transition text-green-700 bg-green-50 border-green-700 border border-transparent"
                  key={tagName}
                  href={`/blog`}
                >
                  {tagName}
                </Link>
              )

            if (tagName === 'All')
              return (
                <Link
                  className="text-sm bg-gray-100 px-3 py-1 rounded-full text-gray-600 transition hover:text-green-700 hover:bg-green-50 hover:border-green-700 border border-transparent"
                  key={tagName}
                  href={`/blog`}
                >
                  {tagName}
                </Link>
              )

            return (
              <Link
                className="text-sm bg-gray-100 px-3 py-1 rounded-full text-gray-600 transition hover:text-green-700 hover:bg-green-50 hover:border-green-700 border border-transparent"
                key={tagName}
                href={`/blog?tag=${tagName}`}
              >
                {tagName}
              </Link>
            )
          })}
      </nav>
      {/* Posts grid */}
      <div className="grid grid-cols-1 gap-4">
        {posts.map((post, index) => {
          const coverImageUrl = (() => {
            if (post.cover && 'external' in post.cover) return post.cover.external.url
            if (post.cover && 'file' in post.cover) return post.cover.file.url
            return null
          })()

          const link = `/blog/${calculateBlogPostSlug(
            post.properties.Title.type === 'title' ? post.properties.Title.title[0].plain_text : '',
          )}`

          return (
            <div
              key={post.id}
              className={`${styles.post} ${index === 0 ? styles.featuredPost : ''}`}
            >
              <div className="p-5 grow grid grid-cols-1 md:grid-cols-6 gap-4">
                <div className="md:col-span-3 order-1 md:order-0 flex flex-col gap-2">
                  <h2 className="text-xl mb-2 transition hover:text-green-700">
                    <Link href={link}>
                      {post.properties.Title.type === 'title' &&
                        post.properties.Title.title[0].plain_text}
                    </Link>
                  </h2>
                  <p className="mb-2">
                    {post.properties.Overview.type === 'rich_text' &&
                      post.properties.Overview.rich_text[0].plain_text}
                  </p>
                  <Link className="transition text-gray-800 hover:text-green-700" href={link}>
                    Read more
                  </Link>
                  <div className="grow" />
                  <div className="flex items-center">
                    <div className="rounded-full shadow overflow-hidden">
                      <Image
                        src={
                          post.properties.Author.type === 'people' &&
                          (post.properties.Author.people[0] as UserObjectResponse).avatar_url
                            ? ((post.properties.Author.people[0] as UserObjectResponse)
                                .avatar_url ?? '/img/logos/app-icon-primary.png')
                            : '/img/logos/app-icon-primary.png'
                        }
                        alt="author's thumbnail"
                        height={50}
                        width={50}
                      />
                    </div>
                    <div className="ml-2 flex flex-col">
                      <p className="text-xs text-gray-700">
                        <strong>
                          {post.properties.Author.type === 'people' &&
                            (post.properties.Author.people[0] as UserObjectResponse).name}
                        </strong>
                      </p>
                      <p className="text-xs">
                        {dayjs(
                          post.properties['Publish Date'].type === 'date'
                            ? post.properties['Publish Date'].date?.start
                            : null,
                        ).format('MMM DD, YYYY')}
                      </p>
                    </div>
                  </div>
                </div>
                <div className="md:col-span-3 order-0 md:order-1">
                  <Link href={link}>
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      className={styles.postCover}
                      src={coverImageUrl ?? ''}
                      alt="article cover image"
                    />
                  </Link>
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default BlogPage
