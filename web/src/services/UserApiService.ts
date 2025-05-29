import { APISettings } from '@/api/config'
import { Client } from '@/api/core/client'
import {
  AppError,
  FormValidationError,
  Result
} from '@/errors/TryCatch'
import { Role } from '@/models/User'
// services/UserApiService.ts

export class UserApiService {
  private client: Client

  constructor () {
    this.client = new Client(APISettings.protocol + APISettings.baseURL, '')
  }

  async loginUser (
    name: string,
    code: string
  ): Promise<
    Result<{
      Token: string
      Role: Role
    }>
  > {
    const resp = await fetch(this.client.url(`/users`), {
      method: 'POST',
      body: JSON.stringify({
        name: name,
        code: code
      })
    })
    if (resp.status === 409) {
      return {
        data: null,
        error: Error(`Пользователь с именем ${name} уже существует!`)
      }
    }

    const json = await resp.json()
    if (!resp.ok) {
      // Проверяем есть ли ошибки формы
      if (json['form_errors']) {
        // Возвращаем ошибки формы в компонент
        return {
          data: null,
          error: new FormValidationError(json['form_errors'])
        }
      }
      throw new Error('Ошибка: В ответе сервера form_errors не указано')
    }

    return { data: { Role: json['role'], Token: json['token'] }, error: null }
  }

  async patchUserRole (username: string, newRole: Role): Promise<void> {
    const resp = await fetch(this.client.url(`/users/${encodeURI(username)}`), {
      method: 'POST',
      // headers: client.headers(),
      body: JSON.stringify({
        Role: newRole
      })
    })

    await resp.json()
    if (!resp.ok) {
      throw new AppError(`статус ответа: ${resp.status}`)
    }
  }

  async removeUser (_: string): Promise<void> {
    // Изолированная логика удаления
  }

  async checkUserToken (token: string): Promise<{
    username: string
    room: string
    role: Role
  }> {
    const resp = await this.client.get('/users/' + token)
    // console.log(resp.status)
    if (resp.status !== 200) {
      throw new AppError('checkUserToken error:')
    }
    const json = await resp.json()
    return {
      username: json['username'],
      role: json['role'],
      room: json['room']
    }
  }
}
