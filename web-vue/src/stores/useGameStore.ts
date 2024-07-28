import { defineStore } from 'pinia'
import { APISettings } from '@/api/config'
import { ref } from 'vue'
import { Client } from '@/api/core/client'
import { useToasterStore } from "../stores/useToasterStore";
import { useUserStore } from "../stores/useUserStore";




export const useGameStore = defineStore('game', () => {
    // actions:
    // {
    const fetching = ref(false)
    const connected = ref(false)
    const roomList = ref<Array<GameInfo>>([])

    // const toasterStore = useToasterStore();
    
    // state: () => {
    return {
        roomList,
        fetching,
        connected,
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