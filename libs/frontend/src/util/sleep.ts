export default async function sleep(ms: number): Promise<void> {
  return await new Promise<void>((res) => setTimeout(() => res(), ms))
}
