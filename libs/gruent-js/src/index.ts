import { Model, Organization, Space, Thing } from '@prisma/client'
import { ListRequestOptions } from './request'
import type { ChannelConfig } from '@gruent/libs/engine/channels/index'
import healthcheck from './healthcheck'
import { getUserById, userSearch } from './users'
import { createSpace, updateSpaceById, getSpaceById, deleteSpaceById } from './spaces'
import { getPolicies, listPolicies, Policy, updatePolicy } from './policies'
import { CreateChannelInput } from './channels/create'
import {
  getOrganizationById,
  updateOrganization,
  createOrganization,
  deleteOrganizationById,
} from './organizations'
import {
  listMembers,
  addMember,
  removeMember,
  changeMemberRole,
  inviteMember,
} from './organizations/members'
import {
  listPermissionSets,
  addPermissionSet,
  updatePermissionSet,
  deletePermissionSet,
} from './permissions'
import { listApiKeys, removeApiKey, addApiKey } from './apiKeys'
import { getChannelExecution, listChannelExecutions } from './channels/executions'
import {
  deleteModel,
  getModel,
  updateModel,
  createModel,
  CreateModelInput,
  listModels,
} from './models'
import {
  channelStatus,
  deleteChannel,
  validateChannel,
  listChannels,
  publishChannel,
  unpublishChannel,
  updateChannelConfig,
  getChannel,
  createChannel,
} from './channels'
import {
  createThing,
  deleteThing,
  listThings,
  getThing,
  healthcheckThing,
  updateThing,
} from './things'
import { getCertificate } from './certificates'
import { listSecrets, createSecret, updateSecret, deleteSecret } from './secrets'
import { listTags, createTag, deleteTag } from './tags'
import applyTagById from './tags/apply'
import removeTagById from './tags/remove'
import {
  signinCredentials,
  signinMagicLink,
  signup,
  verifyMagicLink,
  refresh,
  updatePassword,
} from './auth'

export type ClientOptions = Partial<{
  url: string | URL
}>

export type ClientConfig = {
  url: URL
}

export class Client {
  protected key: string
  private _config: ClientConfig

  constructor(key: string, options: ClientConfig) {
    this.key = key
    this._config = {
      url: new URL(options.url),
    }
  }

  public healthcheck = () => healthcheck(this._config)

  public get auth() {
    return {
      signin: {
        credentials: (email: string, password: string, options?: { redirectTo?: string }) =>
          signinCredentials(this.key, this._config, email, password, options),
        magicLink: {
          send: (email: string, appUrl: string, options?: { redirectTo?: string }) =>
            signinMagicLink(this.key, this._config, email, appUrl, options),
          verify: (token: string) => verifyMagicLink(this.key, this._config, token),
        },
      },
      signup: (email: string, password: string, origin: string, inviteToken?: string) =>
        signup(this.key, this._config, email, password, origin, inviteToken),
      refresh: (refreshToken: string) => refresh(this.key, this._config, refreshToken),
      updatePassword: (newPassword: string) => updatePassword(this.key, this._config, newPassword),
    }
  }

  public get users() {
    return {
      get: (userId: string) => getUserById(this.key, this._config, userId),
      search: (email: string, orgId?: string) => userSearch(this.key, this._config, email, orgId),
    }
  }

  public get organizations() {
    return {
      get: (orgId: string) => getOrganizationById(this.key, this._config, orgId),
      create: (userId: string, name: string) =>
        createOrganization(this.key, this._config, userId, name),
      update: (organization: Organization) =>
        updateOrganization(this.key, this._config, organization),
      delete: (orgId: string) => deleteOrganizationById(this.key, this._config, orgId),
      members: {
        list: (orgId: string, options?: ListRequestOptions) =>
          listMembers(this.key, this._config, orgId, options),
        add: (orgId: string, userId: string) => addMember(this.key, this._config, orgId, userId),
        invite: (orgId: string, email: string, appUrl: string) =>
          inviteMember(this.key, this._config, orgId, email, appUrl),
        remove: (orgId: string, userId: string) =>
          removeMember(this.key, this._config, orgId, userId),
        changeRole: (orgId: string, userId: string, role: 'ADMIN' | 'MEMBER') =>
          changeMemberRole(this.key, this._config, orgId, userId, role),
      },
    }
  }

  public get permissions() {
    return {
      list: (scopeId: string, options?: ListRequestOptions) =>
        listPermissionSets(this.key, this._config, scopeId, options),
      add: (orgId: string, name: string, permissions: string[], spaceId?: string) =>
        addPermissionSet(this.key, this._config, orgId, name, permissions, spaceId),
      update: (
        orgId: string,
        permissionSetId: string,
        name: string,
        permissions: string[],
        spaceId?: string,
      ) =>
        updatePermissionSet(
          this.key,
          this._config,
          orgId,
          permissionSetId,
          name,
          permissions,
          spaceId,
        ),
      remove: (scopeId: string, permissionSetId: string, spaceId?: string) =>
        deletePermissionSet(this.key, this._config, scopeId, permissionSetId, spaceId),
    }
  }

  public get apiKeys() {
    return {
      list: (orgId: string, options?: ListRequestOptions, spaceId?: string) =>
        listApiKeys(this.key, this._config, orgId, options, spaceId),
      add: (orgId: string, permissionSetId: string, name?: string, spaceId?: string) =>
        addApiKey(this.key, this._config, orgId, permissionSetId, name, spaceId),
      remove: (orgId: string, apiKeyId: string, spaceId?: string) =>
        removeApiKey(this.key, this._config, orgId, apiKeyId, spaceId),
    }
  }

  public get spaces() {
    return {
      create: (orgId: string, name: string) => createSpace(this.key, this._config, orgId, name),
      get: (orgId: string, spaceId: string) => getSpaceById(this.key, this._config, orgId, spaceId),
      update: (orgId: string, spaceId: string, space: Space) =>
        updateSpaceById(this.key, this._config, orgId, spaceId, space),
      delete: (orgId: string, spaceId: string) =>
        deleteSpaceById(this.key, this._config, orgId, spaceId),
    }
  }

  public get models() {
    return {
      list: (spaceId: string, options?: ListRequestOptions & { tagFilter?: string[] }) =>
        listModels(this.key, this._config, spaceId, options),
      create: (spaceId: string, modelInput: CreateModelInput) =>
        createModel(this.key, this._config, spaceId, modelInput),
      get: (spaceId: string, modelId: string) => getModel(this.key, this._config, spaceId, modelId),
      update: (spaceId: string, model: Model) =>
        updateModel(this.key, this._config, spaceId, model),
      delete: (spaceId: string, modelId: string) =>
        deleteModel(this.key, this._config, spaceId, modelId),
    }
  }

  public get channels() {
    return {
      list: (spaceId: string, options?: ListRequestOptions & { tagFilter?: string[] }) =>
        listChannels(this.key, this._config, spaceId, options),
      create: (spaceId: string, channelInput: CreateChannelInput) =>
        createChannel(this.key, this._config, spaceId, channelInput),
      get: (spaceId: string, channelId: string) =>
        getChannel(this.key, this._config, spaceId, channelId),
      updateConfig: (spaceId: string, channelId: string, channelConfig: ChannelConfig) =>
        updateChannelConfig(this.key, this._config, spaceId, channelId, channelConfig),
      publish: (spaceId: string, channelId: string) =>
        publishChannel(this.key, this._config, spaceId, channelId),
      unpublish: (spaceId: string, channelId: string) =>
        unpublishChannel(this.key, this._config, spaceId, channelId),
      status: (spaceId: string, channelId: string) =>
        channelStatus(this.key, this._config, spaceId, channelId),
      validate: (spaceId: string, channelConfig: ChannelConfig) =>
        validateChannel(this.key, this._config, spaceId, channelConfig),
      delete: (spaceId: string, channelId: string) =>
        deleteChannel(this.key, this._config, spaceId, channelId),
      executions: {
        list: (
          spaceId: string,
          channelId: string,
          options?: ListRequestOptions & { statusFilter?: string },
        ) => listChannelExecutions(this.key, this._config, spaceId, channelId, options),
        get: (spaceId: string, channelExecutionId: string) =>
          getChannelExecution(this.key, this._config, spaceId, channelExecutionId),
      },
    }
  }

  public get things() {
    return {
      list: (
        spaceId: string,
        options?: ListRequestOptions & {
          category?: string
          filter?: string
          tagFilter?: string[]
        },
      ) => listThings(this.key, this._config, spaceId, options),
      get: (spaceId: string, thingId: string) => getThing(this.key, this._config, spaceId, thingId),
      delete: (spaceId: string, thingId: string) =>
        deleteThing(this.key, this._config, spaceId, thingId),
      create: (spaceId: string, name: string, category: string, attributes: Object) =>
        createThing(this.key, this._config, spaceId, name, category, attributes),
      update: (spaceId: string, thingId: string, thingInput: Pick<Thing, 'name' | 'attributes'>) =>
        updateThing(this.key, this._config, spaceId, thingId, thingInput),
      healthcheck: (spaceId: string, thingCategory: string, attributes: Record<string, any>) =>
        healthcheckThing(this.key, this._config, spaceId, thingCategory, attributes),
    }
  }

  public get policies() {
    return {
      list: (spaceId: string, options?: ListRequestOptions) =>
        listPolicies(this.key, this._config, spaceId, options),
      get: (spaceId: string, certificateId: string) =>
        getPolicies(this.key, this._config, spaceId, certificateId),
      update: (spaceId: string, certificateId: string, policy: Policy, revoke: boolean) =>
        updatePolicy(this.key, this._config, spaceId, certificateId, policy, revoke),
    }
  }

  public get certificates() {
    return {
      get: (spaceId: string, certificateId: string) =>
        getCertificate(this.key, this._config, spaceId, certificateId),
    }
  }

  public get secrets() {
    return {
      list: (spaceId: string) => listSecrets(this.key, this._config, spaceId),
      create: (spaceId: string, name: string, value: string) =>
        createSecret(this.key, this._config, spaceId, name, value),
      update: (spaceId: string, name: string, value: string) =>
        updateSecret(this.key, this._config, spaceId, name, value),
      delete: (spaceId: string, name: string) =>
        deleteSecret(this.key, this._config, spaceId, name),
    }
  }

  public get tags() {
    return {
      list: (spaceId: string, category?: 'things' | 'channels' | 'models') =>
        listTags(this.key, this._config, spaceId, category),
      create: (spaceId: string, name: string) => createTag(this.key, this._config, spaceId, name),
      delete: (spaceId: string, tagId: string) => deleteTag(this.key, this._config, spaceId, tagId),
      apply: (spaceId: string, tagId: string, subjectId: string) =>
        applyTagById(this.key, this._config, spaceId, tagId, subjectId),
      remove: (spaceId: string, tagId: string, subjectId: string) =>
        removeTagById(this.key, this._config, spaceId, tagId, subjectId),
    }
  }
}

export const createClient = (key: string, options?: Partial<ClientOptions>) => {
  // Create default config
  const defaultConfig: ClientConfig = {
    url: new URL('https://api.gruent.com'),
  }

  // Replace defaults with provided options
  const config: ClientConfig = { ...defaultConfig }
  if (options?.url) config.url = new URL(options.url)

  // Create client
  const client = new Client(key, config)

  return client
}
