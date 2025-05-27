import { APISettings } from '@/api/config'
import { Client } from '@/api/core/client'
import { FormValidationError } from '@/errors/TryCatch'
import { Role, UserCreds, UserEntity, UserInfo } from '@/models/User'
import { UserApiService } from '@/services/UserApiService'
import { UserStorageService } from '@/services/UserStorageService'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { useToasterStore } from '../stores/useToasterStore'

UserStorageService.load()

export const useUserStore = defineStore('users', {
  state: () => {
    return {
      // изменение UserCreds вызывает обновление вебсокета
      UserCreds: UserStorageService.load(),
      UserInfo: { room: '', role: Role.Player } as UserInfo,
      fetching: false
    }
  },
  actions: {
    getUser (): UserEntity | null {
      if (!this.UserCreds) return null

      return new UserEntity(
        this.UserCreds.username,
        this.UserCreds.token,
        this.UserInfo.room,
        this.UserInfo.role
      )
    },

    async PatchUser (usr: UserEntity) {
      const toasterStore = useToasterStore()
      const client = new Client(APISettings.protocol + APISettings.baseURL, '')
      let role = ''
      if (usr.role == Role.Player) {
        role = Role.Judge
      } else if (usr.role == Role.Judge) {
        role = Role.Player
      }
      const resp = await fetch(
        client.url(`/users/${encodeURI(usr.username)}`),
        {
          method: 'POST',
          // headers: client.headers(),
          body: JSON.stringify({
            Role: role
          })
        }
      )

      const json = await resp.json()
      if (!resp.ok) {
        console.error('Failed to login with user')
        toasterStore.error(
          `Не удалось изменить роль пользователя ${usr.username}`
        )
        toasterStore.error(json['error'])
        return
      }
    },

    async Remove (usr: string) {
      const toasterStore = useToasterStore()
      const client = new Client(APISettings.protocol + APISettings.baseURL, '')

      const resp = await fetch(client.url(`/users/${encodeURI(usr)}`), {
        method: 'DELETE'
      })
      const json = await resp.json()
      if (!resp.ok) {
        console.error('Failed to login with user')
        toasterStore.error(`Не удалось удалить пользователя ${usr}`)
        toasterStore.error(json['error'])
        return
      }
    },
    async Login (
      input: string,
      code: string
    ): Promise<{
      success: boolean
      formErrors?: Record<string, string>
    }> {
      const userApiService = new UserApiService()
      const toasterStore = useToasterStore()

      const { data, error } = await userApiService.loginUser(input, code)
      if (error) {
        if (error instanceof FormValidationError) {
          return {success: false, formErrors: error.fields}
        }
        return { success: false }
      }

      this.UserCreds = { username: input, token: data.Token }
      this.UserInfo = { room: '', role: data.Role }

      UserStorageService.save(this.UserCreds)

      toasterStore.info(`Вы вошли под именем ${this.UserCreds.username}`)
      console.log(
        `Token for admin ${this.UserCreds.username} created: ${this.UserCreds.token}`
      )
      return { success: true }
    },
    async check () {
      let self = this
      this.fetching = true
      const ok = await (async function () {
        // do something

        if (!self.UserCreds) {
          return false
        }

        const client = new Client(
          APISettings.protocol + APISettings.baseURL,
          ''
        )
        try {
          const token = self.UserCreds.token
          const resp = await client.get('/users/' + token)
          // console.log(resp.status)
          if (resp.status == 200) {
            const json = await resp.json()
            // json["token"] = token
            self.UserCreds = { username: json['username'], token: token }
            self.UserInfo = { room: json['room'], role: json['role'] }
            return true
          } else {
            self.UserCreds = null
          }
        } catch (e) {
          self.UserCreds = null
        }
        return false
      })()
      this.fetching = false
      return ok
    },
    $reset () {
      UserStorageService.clear()
      this.UserCreds = null
    }
  },
  getters: {
    connected: state => state.UserInfo?.room != '',
    hasPermision: state => state.UserInfo.role != Role.Player
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
}
