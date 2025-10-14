export const calculateBlogPostSlug = (title: string): string => {
  return title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

export const slugToNotionQuery = (slug: string): string => {
  return slug.replace(/-/g, ' ')
}

export const calculateHeadingSlug = (heading: string): string => {
  return heading
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}
