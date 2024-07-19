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
    const roomList = ref<Array<RoomInfo>>([])

    const toasterStore = useToasterStore();
    const userStore = useUserStore();
    Fetch()
    async function Fetch() {
        fetching.value = true
        if (userStore.UserCreds == null) {
            return
        }
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        // const token = ref(localStorage.getItem("token") ?? "");
        try {
            const resp = await client.get("/rooms");
            roomList.value = Object.values(await resp.json())
        } catch (e) {
            toasterStore.error("Не удалось обновить информацию о доступных играх")
        }
        fetching.value = false;
    }
    async function CreateGame(roomname: string) {
        if (userStore.UserCreds == null) {
            return;
        }
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        const resp = await fetch(client.url("/rooms"), {
            method: "POST",
            headers: client.headers(),
            body: new URLSearchParams({
                name: roomname,
            }),
        });

        if (!resp.ok) {
            toasterStore.error("не удалось создать игру!");
            return;
        }

        toasterStore.info("Новая игра успешно создана!");
        await Fetch();
    }
    // state: () => {
    return {
        roomList,
        fetching,
        CreateGame,
        Fetch,
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