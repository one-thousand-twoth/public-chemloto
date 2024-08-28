import { useGameStore } from "@/stores/useGameStore"
import { WEBSOCKET_EVENT } from "./websocket"
import { useToasterStore } from "@/stores/useToasterStore"

interface HUB_SUBSCRIBE_EVENT {
    Target: string,
    Name: string
}
export function Subscribe(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("changing to room")
    const b = e.Body as HUB_SUBSCRIBE_EVENT
    if (b.Target == "room"){
        console.log("Hi")
        store.connected = true
        store.name = b.Name
    }
}
export function UNSubscribe(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("exiting room")
    const b = e.Body as HUB_SUBSCRIBE_EVENT
    if (b.Target == "room"){
        console.log("Hi")
        store.connected = false
        store.name = ""
    }
}

export function EngineAction(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("changing to room")
    switch (e.Body["Action"]){
        case "GetElement":{
            console.log(e.Body["Element"])
            store.gameState.Bag.LastElements = e.Body["LastElements"]; 
            store.gameState.Players.forEach((pl) => {pl.Bag[store.currElement] = (pl.Bag[store.currElement] || 0) + 1;} )
            break;
        }
        case "RaiseHand":{
            store.gameState.Players = e.Body["Players"]
            break;
        }
        case "NewTimer":{
            store.timer = e.Body["Value"] as number
            break;
        }
        default:
            console.log("Unresolved EngineAction", e.Body["Action"])

    }

}

export function StartGame(_: WEBSOCKET_EVENT){
    const toaster = useToasterStore()
    const store = useGameStore()
    store.gameState.Started = true
    toaster.info(
        "Игра начата!"
    )
}