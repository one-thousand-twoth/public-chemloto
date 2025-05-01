import { WebsocketConnector } from '@/api/websocket/websocket'
import { BaseStateHandler as BaseStateController } from './base'

// Реализация обработчика для состояния OBTAIN
export default class ObtainStateController extends BaseStateController<'OBTAIN'> {
  constructor (ws: WebsocketConnector) {
    super(ws, 'OBTAIN')
  }

  getElement (): void {
    this.sendCommand({ Action: 'GetElement', Type: 'ENGINE_ACTION' })
  }

  sendContinue (): void {
    this.sendCommand({ Action: 'Continue', Type: 'ENGINE_ACTION' })
  }
}
