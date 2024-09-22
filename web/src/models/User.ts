export enum Role {
    Admin = "Admin_Role",
    Judge = "Judge_Role",
    Player = "Player_Role",
}
export function i18nRole(role: Role) {
    switch (role) {
        case Role.Admin:
            return "Админ"
        case Role.Judge:
            return "Судья"
        case Role.Player:
            return "Player"
    }
}
export function emojiRole(role: Role) {
    switch (role) {
        case Role.Admin:
            return "🅰"
        case Role.Judge:
            return "🛠"
        case Role.Player:
            return "✌"
    }
}
export interface UserInfo {
    username: string
    token: string
    role: Role
    room: string
    // status: string
}