import isBrowser from './isBrowser'

const copyToClipboard = async (text: string) => {
  if (!isBrowser) return { error: null }

  try {
    await navigator.clipboard.writeText(text)
  } catch (e) {
    return { error: e }
  }

  return { error: null }
}

export default copyToClipboard
