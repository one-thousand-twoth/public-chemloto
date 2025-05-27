import { useGameStore } from "@/stores/useGameStore";
import { useToasterStore } from "@/stores/useToasterStore";
import { EngineAction, StartGame, Subscribe, UNSubscribe as UnSubscribe } from "./handlers";

export interface WEBSOCKET_EVENT {
    Type: string
    Ok: boolean
    Errors: Array<string>
    Body: { [id: string]: any; }
}

export class WebsocketConnector {
    baseUrl: string;
    token: string;
    active: boolean;
    connection!: WebSocket | null;
    // gameStore: Store;
    // gameStore: 

    constructor(baseUrl: string, token: string = "") {
        this.token = token;
        this.baseUrl = baseUrl;
        this.active = false;
    }
    Run() {
        const gameStore = useGameStore()
        const toaster = useToasterStore()

        this.connection = new WebSocket(`ws://${this.baseUrl}/api/v2/ws?token=${this.token}`)
        this.connection.onmessage = function (event) {
            const data = JSON.parse(event.data) as WEBSOCKET_EVENT
            if (!data.Ok) {
                data.Errors.forEach((err) => {
                    toaster.error(err)
                })
                return
            }
            switch (data.Type) {
                case "HUB_SUBSCRIBE":
                    Subscribe(data)
                    break;
                case "HUB_UNSUBSCRIBE":
                    UnSubscribe(data)
                    break;
                case "HUB_STARTGAME":
                    StartGame(data)
                    break;

                case "ENGINE_ACTION":
                    EngineAction(data)
                    break;
                case "ENGINE_INFO":
                    gameStore.EngineInfo(data)

            }
            console.log(data)
        }

        this.connection.onopen = function (event) {
            console.log(event)
            console.log("Successfully connected to the websocket server...")
        }
    }
    Send(msg: Object) {
        if (!this.connection) {
            console.error("error sending message to websocket, no connecttion")
            return
        }
        this.connection.send(JSON.stringify(msg))
    }
}