import { ChannelEdge, ChannelNote, ChannelRuntime } from '@gruent/libs/engine/channels/channels'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { create } from 'zustand'
import { IoPoint } from '../components/organisms/ChannelEditor/types/IoPoint'
import { Model, Thing } from '@prisma/client'

type State = {
  nodes: Map<string, ChannelNode>
  ioPoints: Map<string, IoPoint>
  edges: Map<string, ChannelEdge>
  notes: Map<string, ChannelNote>
  channelName: string
  currentNode: ChannelNode | null
  currentEdge: ChannelEdge | null
  currentNote: ChannelNote | null
  runtime: ChannelRuntime
  validationErrors: {
    channelErrors: string[] | null
    nodeErrors: Record<string, string[]> | null
  } | null
  pan: { x: number; y: number }
  status?: 'PUBLISHED' | 'UNPUBLISHED' | 'INITIALIZING'
  zoom: number
  edgePreview: {
    from: { x: number; y: number }
    to: { x: number; y: number }
  } | null
  models: Model[]
  things: Thing[]
}

type Action = {
  reset: () => void
  upsertNode: (nodeId: string, node: ChannelNode) => void
  removeNode: (nodeId: string) => void
  upsertIoPoint: (ioPointId: string, ioPoint: IoPoint) => void
  removeIoPoint: (ioPointId: string) => void
  upsertEdge: (edgeId: string, edge: ChannelEdge) => void
  removeEdge: (edgeId: string) => void
  upsertNote: (noteId: string, note: ChannelNote) => void
  removeNote: (noteId: string) => void
  setChannelName: (channelName: string) => void
  setCurrentNode: (currentNode: ChannelNode | null) => void
  setCurrentEdge: (edge: ChannelEdge | null) => void
  setCurrentNote: (note: ChannelNote | null) => void
  setRuntime: (runtime: ChannelRuntime) => void
  setValidationErrors: (
    validationErrors: {
      channelErrors: string[] | null
      nodeErrors: Record<string, string[]> | null
    } | null,
  ) => void
  setPan: (x: number, y: number) => void
  setStatus: (status: 'PUBLISHED' | 'UNPUBLISHED' | 'INITIALIZING') => void
  setZoom: (zoom: number) => void
  setEdgePreview: (
    edgePreview: { from: { x: number; y: number }; to: { x: number; y: number } } | null,
  ) => void
  setModels: (models: Model[]) => void
  setThings: (things: Thing[]) => void
}

const initialState: State = {
  nodes: new Map([]),
  edges: new Map([]),
  ioPoints: new Map([]),
  notes: new Map([]),
  channelName: '',
  currentNode: null,
  currentEdge: null,
  currentNote: null,
  runtime: {
    size: 'small',
  },
  validationErrors: null,
  pan: { x: 0, y: 0 },
  zoom: 1,
  edgePreview: null,
  models: [],
  things: [],
}

const useChannelEditorStore = create<State & Action>((set, get) => ({
  ...initialState,
  reset: () => set(initialState),
  upsertNode: (nodeId, node) => {
    const nodes = get().nodes
    const updatedNodes = new Map(nodes)
    updatedNodes.delete(nodeId) // deleting plus setting causes the node to go to the top of the map, rendering first in draw.tsx
    updatedNodes.set(nodeId, node)
    set({ nodes: updatedNodes })
  },
  removeNode: (nodeId) => {
    const nodes = get().nodes
    const updatedNodes = new Map(nodes)
    updatedNodes.delete(nodeId)
    set({ nodes: updatedNodes })
  },
  upsertIoPoint: (ioPointId, ioPoint) => {
    const ioPoints = get().ioPoints
    const updatedIoPoints = new Map(ioPoints)
    updatedIoPoints.delete(ioPointId) // deleting plus setting causes the ioPoint to go to the top of the map, rendering first in draw.tsx
    updatedIoPoints.set(ioPointId, ioPoint)
    set({ ioPoints: updatedIoPoints })
  },
  removeIoPoint: (ioPointId) => {
    const ioPoints = get().ioPoints
    const updatedIoPoints = new Map(ioPoints)
    updatedIoPoints.delete(ioPointId)
    set({ ioPoints: updatedIoPoints })
  },
  upsertEdge: (edgeId, edge) => {
    const edges = get().edges
    const updatedEdges = new Map(edges)
    updatedEdges.delete(edgeId) // deleting plus setting causes the edge to go to the top of the map, rendering first in draw.tsx
    updatedEdges.set(edgeId, edge)
    set({ edges: updatedEdges })
  },
  removeEdge: (edgeId) => {
    const edges = get().edges
    const updatedEdges = new Map(edges)
    updatedEdges.delete(edgeId)
    set({ edges: updatedEdges })
  },
  upsertNote: (noteId, note) => {
    const notes = get().notes
    const updatedNotes = new Map(notes)
    updatedNotes.delete(noteId) // deleting plus setting causes the note to go to the top of the map, rendering first in draw.tsx
    updatedNotes.set(noteId, note)
    set({ notes: updatedNotes })
  },
  removeNote: (noteId) => {
    const notes = get().notes
    const updatedNotes = new Map(notes)
    updatedNotes.delete(noteId)
    set({ notes: updatedNotes })
  },
  setChannelName: (channelName) => set({ channelName }),
  setCurrentNode: (currentNode) => set({ currentNode, currentEdge: null, currentNote: null }),
  setCurrentEdge: (currentEdge) => set({ currentNode: null, currentEdge, currentNote: null }),
  setCurrentNote: (currentNote) => set({ currentNode: null, currentEdge: null, currentNote }),
  setRuntime: (runtime) => set({ runtime }),
  setValidationErrors: (validationErrors) => set({ validationErrors }),
  setPan: (x: number, y: number) => set({ pan: { x, y } }),
  setStatus: (status) => set({ status }),
  setZoom: (zoom) => set({ zoom }),
  setEdgePreview: (edgePreview) => set({ edgePreview }),
  setModels: (models) => set({ models }),
  setThings: (things) => set({ things }),
}))

export default useChannelEditorStore
