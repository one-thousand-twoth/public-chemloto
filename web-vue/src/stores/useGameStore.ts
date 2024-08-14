import { defineStore } from 'pinia'
import { ref } from 'vue'




export const useGameStore = defineStore('game', () => {
    // actions:
    // {
    const fetching = ref(false)
    const connected = ref(false)
    const roomList = ref<Array<GameInfo>>([])
    const LastElements = ref<Array<string>>(Array.from({ length: 5 }, () => "UNDEFINED"))
    const currElement = ref('UNDEFINED')

    // const toasterStore = useToasterStore();
    
    // state: () => {
    return {
        roomList,
        fetching,
        connected,
        currElement,
        LastElements
    }
    // },
    // Add(r: RoomInfo) {
    //     this.RoomList.push(r)
    // },
    // AddAll(r: RoomInfo[]) {
    //     this.RoomList.push(...r)
    // },
})

export interface GameInfo {
    name: string
    // status: string
} 