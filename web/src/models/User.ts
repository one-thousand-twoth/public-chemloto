import { ExclamationCircleIcon, EyeIcon, StarIcon, UserIcon } from "@heroicons/vue/16/solid";
import { FunctionalComponent, HTMLAttributes, VNodeProps } from "vue";
export enum Role {
    Admin = "Admin_Role",
    Judge = "Judge_Role",
    Player = "Player_Role",
}

const roleData = {
    [Role.Admin]: { label: "Админ", emoji: "🅰", icon: EyeIcon },
    [Role.Judge]: { label: "Судья", emoji: "🛠", icon: StarIcon },
    [Role.Player]: { label: "Игрок", emoji: "✌", icon: UserIcon },
}


export function i18nRole(role: Role) {
    return roleData[role]?.label || "(?)"
}
export function IconRole(role: Role): FunctionalComponent<HTMLAttributes & VNodeProps> {
    return roleData[role]?.icon || ExclamationCircleIcon
}
export type User = UserCreds & UserInfo

export type UserCreds =  {
    username: string
    token: string

    // room: string
    // status: string
}

export type UserInfo = {
    role: Role
    room: string
}