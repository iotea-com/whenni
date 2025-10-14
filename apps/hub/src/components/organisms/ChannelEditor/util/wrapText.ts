export const wrapText = (context: CanvasRenderingContext2D, text: string, maxWidth: number) => {
  const words = text.split(' ')
  let line = ''
  const lines: string[] = []

  for (let i = 0; i < words.length; i++) {
    const testLine = line + words[i] + ' '
    const metrics = context.measureText(testLine)
    const testWidth = metrics.width

    if (testWidth > maxWidth && i > 0) {
      lines.push(line)
      line = words[i] + ' '
      continue
    } else {
      line = testLine
    }
  }

  lines.push(line)

  return lines
}
