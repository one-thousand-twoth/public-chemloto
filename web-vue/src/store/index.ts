import { defineStore } from 'pinia'

export const useRoomsStore = defineStore('rooms', {
    state: () => {
        return {
            // for initially empty lists
            RoomList: [
            {
                name: "Test",
                status: "ok"
            }
            ] as RoomInfo[],
        }
    },
    actions: {
        Add(r: RoomInfo) {
            this.RoomList.push(r)
        },
        AddAll(r: RoomInfo[]) {
            this.RoomList.push(...r)
        },
    }
})

interface RoomInfo {
    name: string
    status: string
}