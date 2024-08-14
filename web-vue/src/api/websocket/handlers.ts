import { useGameStore } from "@/stores/useGameStore"
import { WEBSOCKET_EVENT } from "./websocket"

interface SUBSCRIBE_EVENT {
    Target: string,
    Name: string
}
export function Subscribe(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("changing to room")
    const b = e.Body as SUBSCRIBE_EVENT
    if (b.Target == "room"){
        console.log("Hi")
        store.connected = true
    }
}

export function EngineAction(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("changing to room")
    switch (e.Body["Action"]){
        case "GetElement":{
            console.log(e.Body["Element"])
            store.currElement = e.Body["Element"]
            store.LastElements = e.Body["LastElements"].reverse().concat(Array.from({ length: 5 - e.Body["LastElements"].length }, () => "UNDEFINED")); 
            break;
        }
        default:
            console.log("Unresolved EngineAction", e.Body["Action"])

    }

}
