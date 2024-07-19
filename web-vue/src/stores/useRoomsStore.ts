import { defineStore } from 'pinia'
import { APISettings } from '@/api/config'
import { ref } from 'vue'
import { Client } from '@/api/core/client'
import { useToasterStore } from "../stores/useToasterStore";
import { useUserStore } from "../stores/useUserStore";




export const useRoomsStore = defineStore('rooms', () => {
    // actions:
    // {
    const fetching = ref(false)
    const RoomList = ref<Array<RoomInfo>>([])

    const toasterStore = useToasterStore();
    const userStore = useUserStore();

    async function Fetch() {
        fetching.value = true
        if (userStore.UserCreds == null) {
            return
        }
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        // const token = ref(localStorage.getItem("token") ?? "");
        try {
            const resp = await client.get("/rooms");
            RoomList.value = Object.values(await resp.json())
        } catch (e) {
            toasterStore.error("Не удалось обновить информацию о доступных играх")
        }
        fetching.value = false
    }
    async function CreateGame() {
        fetching.value = true
        if (userStore.UserCreds == null) {
            return
        }
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        const resp = await fetch(client.url("/rooms"), {
            method: "POST",
            headers: client.headers(),
            body: new URLSearchParams({
                name: "new-room",
                IP: "172.16.1.126",
            }),
        });

        if (!resp.ok) {
            toasterStore.error("не удалось создать игру!");
        }
    
        toasterStore.info("Новая игра успешно создана!");
        await Fetch()
        fetching.value = false
    }
    // state: () => {
    return {
        RoomList: [] as RoomInfo[],
        Fetching: false,
        CreateGame,
        Fetch
    }
    // },
    // Add(r: RoomInfo) {
    //     this.RoomList.push(r)
    // },
    // AddAll(r: RoomInfo[]) {
    //     this.RoomList.push(...r)
    // },
})

interface RoomInfo {
    name: string
    // status: string
} 