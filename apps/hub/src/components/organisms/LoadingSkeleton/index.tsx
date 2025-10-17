const LoadingSkeleton = () => {
  return (
    <div className="w-full max-w-lg">
      <div
        className={`inline-block w-3/4 h-4 bg-gray-200 dark:bg-gray-800 rounded-sm animate-shimmer`}
      />
      <div
        className={`inline-block w-full h-4 bg-gray-200 dark:bg-gray-800 rounded-sm animate-shimmer`}
      />
      <div
        className={`inline-block w-full h-4 bg-gray-200 dark:bg-gray-800 rounded-sm animate-shimmer`}
      />
      <div
        className={`inline-block w-full h-4 bg-gray-200 dark:bg-gray-800 rounded-sm animate-shimmer`}
      />
    </div>
  )
}

export default LoadingSkeleton
