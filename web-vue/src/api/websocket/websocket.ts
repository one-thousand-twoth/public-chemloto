import { Subscribe } from "./handlers";

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
        // const gameStore = useGameStore()
        this.connection = new WebSocket(`ws://${this.baseUrl}/api/v1/ws?token=${this.token}`)
        this.connection.onmessage = function (event) {
            const data = JSON.parse( event.data )
            if (data["Type"] == "HUB_SUBSCRIBE"){
                Subscribe(data)
            }
            console.log(data)
        }

        this.connection.onopen = function (event) {
            console.log(event)
            console.log("Successfully connected to the echo websocket server...")
        }
    }
    Send(msg: Object){
        if (!this.connection){
            return
        }
        this.connection.send(JSON.stringify(msg))
    }
}