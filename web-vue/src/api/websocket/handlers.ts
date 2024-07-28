import { useGameStore } from "@/stores/useGameStore"

export function Subscribe(e: any) {
    const store = useGameStore()
    console.log("changing to room")
    e = e as SUBSCRIBE_EVENT
    if (e.Target == "room"){
        console.log("Hi")
        store.connected = true
    }
}

interface SUBSCRIBE_EVENT {
    Type: string,
    Target: string,
    Name: string
}