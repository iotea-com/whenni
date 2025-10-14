export type ModelAttributes = Record<string, ModelAttribute>

export type ModelAttribute = {
  id: string
  key: string
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  subtype?: 'string' | 'number' | 'boolean' | 'object' | string
  arrayObjectAttributes?: Record<string, ModelAttribute>
  required: boolean
  protected: boolean
  parentId?: string
  order?: number
}
