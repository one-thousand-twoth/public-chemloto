import { APISettings } from '@/api/config';
import { WebsocketConnector } from '@/api/websocket/websocket';
import { WEBSOCKET_EVENT } from '@/models/Messages';
import { defineStore } from 'pinia';

type EventHandler = (ev: WEBSOCKET_EVENT) => void;

export const useWebSocketStore = defineStore('webSocket', {
  state: () => ({
    messages: [] as Array<any>,
    connector:  new WebsocketConnector(APISettings.baseURL, '') as WebsocketConnector,
    listeners: {} as { [key: string]: EventHandler[]},
  }),
  actions: {
    addMessage(message: any) {
      this.messages.push(message);
    },
    initializeWebSocket() {
      websocketService.setOnMessage((data) => {
        this.addMessage(data);
      });
    },
    sendMessage(message: object) {
      websocketService.sendMessage(message);
    },
  },
});