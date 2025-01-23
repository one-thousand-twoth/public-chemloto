import { StartTimer, useGameStore } from "@/stores/useGameStore"
import { useToasterStore } from "@/stores/useToasterStore"
import { useUserStore } from "@/stores/useUserStore"
import { WEBSOCKET_EVENT } from "./websocket"

interface HUB_SUBSCRIBE_EVENT {
    Target: string,
    Name: string
}
export function Subscribe(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    const user = useUserStore()
    console.log("changing to room")
    const b = e.Body as HUB_SUBSCRIBE_EVENT
    if (b.Target == "room") {
        console.log("Hi")
        user.UserInfo!.room = b.Name
        store.name = b.Name
    }
}
export function UNSubscribe(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    const user = useUserStore()
    const b = e.Body as HUB_SUBSCRIBE_EVENT
    if (b.Target == "room") {
        console.log("exiting room", b)
        if (user.UserCreds == null) {
            console.error("user is null")
            return
        }
        user.UserInfo.room = ""
        console.log(user.UserCreds)
        store.name = ""
    }
}

export function EngineAction(e: WEBSOCKET_EVENT) {
    const store = useGameStore()
    console.log("handling engine action")
    switch (e.Body["Action"]) {
        case "GetElement": {
            console.log(e.Body["Element"])
            store.gameState.Bag.LastElements = e.Body["LastElements"];
            store.gameState.Players.forEach((pl) => { pl.Bag[store.currElement] = (pl.Bag[store.currElement] || 0) + 1; })

            break;
        }
        case "RaiseHand": {
            store.gameState.Players = e.Body["Players"]
            break;
        }
        case "NewTimer": {
            store.timer = e.Body["Value"] as number
            StartTimer()
            break;
        }
        default:
            console.log("Unresolved EngineAction", e.Body["Action"])

    }

}

export function StartGame(_: WEBSOCKET_EVENT) {
    const toaster = useToasterStore()
    const store = useGameStore()
    store.gameState.Status = "STATUS_STARTED"
    toaster.info(
        "Игра начата!"
    )
}