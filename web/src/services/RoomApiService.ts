import { APISettings } from '@/api/config'
import { Client } from '@/api/core/client'
import { assertIsErrorResponse } from '@/errors/protocol_errors'
import { AppError, FormValidationError } from '@/errors/TryCatch'
import { CreateRoomRequest } from '@/models/RoomModel'
import { useUserStore } from '@/stores/useUserStore'


export class RoomApiService {
  private client: Client

  constructor () {
    const userStore = useUserStore()
    const token = userStore.UserCreds?.token
    this.client = new Client(
      APISettings.protocol + APISettings.baseURL,
      token ?? ''
    )
  }

  async createRoom (room: CreateRoomRequest): Promise<void> {
    const resp = await fetch(this.client.url('/rooms'), {
      method: 'POST',
      headers: this.client.headers(),
      body: JSON.stringify(room)
    })
    const json = await resp.json()

    if (!resp.ok) {
      assertIsErrorResponse(json)
      if (json.form_errors){
        throw new FormValidationError(json.form_errors)
      }
      throw new AppError(json.error)
    }
  }
}
