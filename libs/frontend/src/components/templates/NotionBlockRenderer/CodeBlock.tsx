import { createStarryNight, common } from '@wooorm/starry-night'
import { toJsxRuntime } from 'hast-util-to-jsx-runtime'
import { Fragment, jsx, jsxs } from 'react/jsx-runtime'

const CodeBlock = async ({ content, language }: { content: string; language: string }) => {
  if (!content) return null

  const starryNight = await createStarryNight(common)
  const scope = starryNight.flagToScope(language)

  if (scope) {
    const tree = starryNight.highlight(content, scope)
    const reactNode = toJsxRuntime(tree, { Fragment, jsx, jsxs })

    return <code>{reactNode}</code>
  }

  return <code>{content}</code>
}

export default CodeBlock
