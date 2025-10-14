export default function Loading() {
  return (
    <div className="flex flex-col gap-4 px-8 pt-4">
      <div className="grid grid-cols-4 gap-4">
        <div className="h-8 bg-gray-200 dark:bg-gray-800 rounded animate-shimmer" />
        <div />
        <div />
        <div className="h-8 bg-gray-200 dark:bg-gray-800 rounded animate-shimmer" />
      </div>
      <div className="grid grid-cols-5 gap-4">
        <div className="col-span-3 h-[400px] bg-gray-200 dark:bg-gray-800 rounded animate-shimmer" />
        <div className="col-span-2 h-[400px] bg-gray-200 dark:bg-gray-800 rounded animate-shimmer" />
      </div>
    </div>
  )
}
