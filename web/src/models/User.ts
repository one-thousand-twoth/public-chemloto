import {
  ExclamationCircleIcon,
  EyeIcon,
  StarIcon,
  UserIcon
} from '@heroicons/vue/16/solid'
import { FunctionalComponent, HTMLAttributes, VNodeProps } from 'vue'
export enum Role {
  Admin = 'Admin_Role',
  Judge = 'Judge_Role',
  Player = 'Player_Role'
}

const roleData = {
  [Role.Admin]: { label: 'Админ', emoji: '🅰', icon: EyeIcon },
  [Role.Judge]: { label: 'Судья', emoji: '🛠', icon: StarIcon },
  [Role.Player]: { label: 'Игрок', emoji: '✌', icon: UserIcon }
}

export function i18nRole (role: Role) {
  return roleData[role]?.label || '(?)'
}
export function IconRole (
  role: Role
): FunctionalComponent<HTMLAttributes & VNodeProps> {
  return roleData[role]?.icon || ExclamationCircleIcon
}
export class UserEntity {
  constructor (
    public username: string,
    public token: string,
    public room: string,
    public role: Role
  ) {}

  hasPermission (): boolean {
    return this.role !== Role.Player
  }
  public switchRole () {
    let role: Role
    if (this.role == Role.Player) {
      role = Role.Judge
    } else {
      role = Role.Player
    }
    this.role = role
  }
  static fromJSON (obj: any): UserEntity {
    return new UserEntity(obj.username, obj.token, obj.room, obj.role)
  }
}

export type UserCreds = {
  username: string
  token: string

  // room: string
  // status: string
}

export type UserInfo = {
  role: Role
  room: string
}
