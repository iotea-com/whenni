export type ApiResponse<T> = {
  data: T | null
  errors: string[] | null
  page?: number
  totalPages?: number
  totalResults?: number
  resultsPerPage?: number
}
