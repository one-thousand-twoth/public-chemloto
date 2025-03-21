import { APISettings } from '@/api/config';
import { Client } from '@/api/core/client';
import { Room } from '@/models/RoomModel';
import { acceptHMRUpdate, defineStore } from 'pinia';
import { ref } from 'vue';
import { useToasterStore } from "../stores/useToasterStore";
import { useUserStore } from "../stores/useUserStore";




export const useRoomsStore = defineStore('rooms', () => {
    const fetching = ref(false)
    const roomList = ref<Array<Room>>([])

    const toasterStore = useToasterStore();
    const userStore = useUserStore();
    Fetch()
    async function Fetch() {
        fetching.value = true
        if (userStore.UserCreds == null) {
            fetching.value = false;
            return
        }
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        // const token = ref(localStorage.getItem("token") ?? "");
        try {
            const resp = await client.get("/rooms");
            if (resp.status == 200) {
                roomList.value = Object.values(await resp.json())
            }
        } catch (e) {
            toasterStore.error("Не удалось обновить информацию о доступных играх")
        }
        fetching.value = false;
    }
    async function CreateGame(room: Room) {
        if (userStore.UserCreds == null) {
            return;
        }
        console.log(room)
        const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
        const resp = await fetch(client.url("/rooms"), {
            method: "POST",
            headers: client.headers(),
            body: JSON.stringify(room),
        });

        if (!resp.ok) {
            const json = await resp.json()
            json['error'].forEach((element: string) => {
                toasterStore.error(element);
            });
            return;
        }

        toasterStore.info("Новая игра успешно создана!");
        await Fetch();
    }
    return {
        roomList,
        fetching,
        CreateGame,
        Fetch,
    }
})

// export interface RoomInfo {
//     name: string
//     maxPlayers: number
//     engine: {Status: string}
//     elementCounts: { [id: string]: number; }
//     time: number
//     isAuto: boolean,
//     isAutoCheck: boolean,
// }

export function i18nStatus(status: string): string {
    switch (status) {
        case "STATUS_WAITING":
            return "Ожидает"
        case "STATUS_STARTED":
            return "Запущена"
        case "STATUS_COMPLETED":
            return "Завершена"
        default:
            return "??";
    }
}

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useRoomsStore, import.meta.hot))
}
