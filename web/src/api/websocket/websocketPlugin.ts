// websocketPlugin.ts
import { Plugin } from 'vue'
import { WebsocketConnector } from '@/api/websocket/websocket'
import { PiniaPluginContext } from 'pinia'

export interface StoreWithWS {
    $ws: WebsocketConnector
}

export const websocketPlugin: Plugin = {
    install: (app, wsConnection: WebsocketConnector) => {
        app.provide('websocket', wsConnection)
    }
}

export const piniaWebsocketPlugin = (wsConnection: WebsocketConnector) => {
    return ({ store }: PiniaPluginContext) => {
        store.$ws = wsConnection
    }
}
