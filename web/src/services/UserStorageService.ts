import { UserCreds } from '@/models/User'

// services/UserStorageService.ts
export class UserStorageService {
  private static readonly USER_KEY = 'user'

  static save (userCreds: UserCreds): void {
    localStorage.setItem(this.USER_KEY, JSON.stringify(userCreds))
  }

  static load (): UserCreds | null {
    const value = localStorage.getItem(this.USER_KEY)
    return value ? JSON.parse(value) : null
  }

  static clear (): void {
    localStorage.removeItem(this.USER_KEY)
  }
}
