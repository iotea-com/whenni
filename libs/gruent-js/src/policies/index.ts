export { default as getPolicies } from './get'
export { default as listPolicies } from './list'
export { default as updatePolicy } from './update'

// TODO: move to DB schema or shared library
export type Policy = {
  allowedSubscriptionTopics: string[]
  allowedPublishTopics: string[]
}
