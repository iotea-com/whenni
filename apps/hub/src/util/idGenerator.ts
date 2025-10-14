import { customAlphabet } from 'nanoid'
const nanoid = customAlphabet('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789')

export function generateAccountId() {
  return `a${nanoid(15)}`
}

export function generateUserId() {
  return `u${nanoid(15)}`
}

export function generateEdgeId() {
  return `e${nanoid()}`
}

export function generateModelAttributeId() {
  return `f${nanoid()}`
}

export function generateNodeId() {
  return `n${nanoid()}`
}

export function generateNoteId() {
  return `note_${nanoid()}`
}
