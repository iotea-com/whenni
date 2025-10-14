const toSlug = (text: string) => {
  return text.toLowerCase().replace(/\W+/g, '-')
}

export default toSlug
