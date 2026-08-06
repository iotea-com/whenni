import { ApiResponse } from './response'

export type ListRequestOptions = {
  page?: number
  resultsPerPage?: number
  filter?: string
}

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD'

type RequestOptions = {
  authorization?: string
  headers?: Record<string, string>
  query?: URLSearchParams
}

async function apiRequest<ResponseData>(
  route: URL,
  method: HttpMethod,
  body: Record<string, any> | null,
  options?: RequestOptions,
): Promise<ApiResponse<ResponseData>> {
  // Init fetch request
  const requestOptions: RequestInit = {
    method,
  }

  // Set request body
  if (body) requestOptions.body = JSON.stringify(body)

  // Set request headers
  const headers = (() => {
    const h: Record<string, string> = options && options.headers ? { ...options.headers } : {}

    if (options && options.authorization) {
      h['Authorization'] = `Bearer ${options.authorization}`
    }

    return h
  })()

  if (Object.keys(headers).length > 0) requestOptions.headers = headers

  // Set search params
  if (options && options.query) {
    options.query.forEach((v, k) => {
      route.searchParams.set(k, v)
    })
  }

  try {
    // Send request
    const response = await fetch(route.toString(), requestOptions)

    // Handle 500 errors
    if (response.status === 500) {
      return {
        data: null,
        errors: [await response.text()],
      }
    }

    // Handle 422 errors
    if (response.status === 422) {
      return {
        data: null,
        errors: ['invalid request body'],
      }
    }

    // Handle 400 errors
    if (response.status === 400) {
      try {
        const responseBody = (await response.json()) as ApiResponse<ResponseData>

        return {
          data: null,
          errors: responseBody.errors ?? ['invalid request'],
        }
      } catch (_e) {
        return {
          data: null,
          errors: ['invalid request'],
        }
      }
    }

    // Handle 401 errors
    if (response.status === 401) {
      const error = await (async () => {
        try {
          const responseBody = (await response.json()) as ApiResponse<ResponseData>
          return responseBody.errors ?? ['Unauthorized request']
        } catch (_e) {
          return ['Unauthorized request']
        }
      })()

      return {
        data: null,
        errors: error,
      }
    }

    // Handle 404 errors
    if (response.status === 404) {
      return {
        data: null,
        errors: ['API route not found'],
      }
    }

    // Handle 200-499 responses from the API
    const responseBody = (await response.json()) as ApiResponse<ResponseData>

    return {
      data: responseBody.data,
      errors: responseBody.errors,
      page: responseBody.page,
      totalPages: responseBody.totalPages,
      totalResults: responseBody.totalResults,
      resultsPerPage: responseBody.resultsPerPage,
    }
  } catch (_e) {
    return {
      data: null,
      errors: ['API is not reachable'],
    }
  }
}

export default apiRequest
