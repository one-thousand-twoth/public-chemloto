import { json } from "stream/consumers";

export class WebsocketConnector {
    baseUrl: string;
    token: string;
    active: boolean;
    connection!: WebSocket | null;

    constructor(baseUrl: string, token: string = "") {
        this.token = token;
        this.baseUrl = baseUrl;
        this.active = false;
    }
    Run() {
        this.connection = new WebSocket(`ws://${this.baseUrl}/api/v1/ws?token=${this.token}`)
        this.connection.onmessage = function (event) {
            console.log(event);
            console.log(event.data)
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