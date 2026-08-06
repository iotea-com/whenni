import { ClientConfig } from '..'
import { ApiResponse } from '../response'

type ResponseData = {
  blob: Blob
  filename: string
}

const getCertificate = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  certificateId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/certificates/${certificateId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  url.search = queryParams.toString()

  const response = await fetch(url, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${key}`,
    },
  })

  if (response.status !== 200) {
    return {
      data: null,
      errors: [await response.text()],
    }
  }

  const blob = await response.blob()

  const filename = (() => {
    const contentDisposition = response.headers.get('Content-Disposition')
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename=\"(.*)\"/)
      if (filenameMatch && filenameMatch.length > 1) return filenameMatch[1]
    }

    return 'certificate_bundle.zip'
  })()

  return {
    data: {
      blob,
      filename,
    },
    errors: null,
  }
}

export default getCertificate
