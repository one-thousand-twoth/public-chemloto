import { APISettings } from '@/api/config'
import { Client } from '@/api/core/client'
import { AppError, FormValidationError } from '@/errors/TryCatch'
import { CreateRoomRequest, Room } from '@/models/RoomModel'
import { RoomApiService } from '@/services/RoomApiService'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { ref } from 'vue'
import { useToasterStore } from '../stores/useToasterStore'
import { useUserStore } from '../stores/useUserStore'

export const useRoomsStore = defineStore('rooms', () => {
  const fetching = ref(false)
  const roomList = ref<Array<Room>>([])

  const toasterStore = useToasterStore()
  const userStore = useUserStore()
  Fetch()
  async function Fetch () {
    fetching.value = true
    if (userStore.UserCreds == null) {
      fetching.value = false
      return
    }
    const client = new Client(
      APISettings.protocol + APISettings.baseURL,
      userStore.UserCreds.token
    )
    // const token = ref(localStorage.getItem("token") ?? "");
    try {
      const resp = await client.get('/rooms')
      if (resp.status == 200) {
        roomList.value = Object.values(await resp.json())
      }
    } catch (e) {
      toasterStore.error('Не удалось обновить информацию о доступных играх')
    }
    fetching.value = false
  }
  async function CreateGame (
    room: CreateRoomRequest
  ): Promise<Record<string, string> | null> {
    if (userStore.UserCreds == null) {
      throw new AppError('Вызов без UserCreds')
    }
    const api = new RoomApiService()
    try {
      await api.createRoom(room)
    } catch (e) {
      if (e instanceof FormValidationError) {
        return e.fields
      }
      if (e instanceof AppError) {
        toasterStore.error(e.message)
        return null
      }
      throw e
    }

    toasterStore.info('Новая игра успешно создана!')
    await Fetch()
    return null
  }
  return {
    roomList,
    fetching,
    CreateGame,
    Fetch
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useRoomsStore, import.meta.hot))
}
