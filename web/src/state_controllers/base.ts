import { WebsocketConnector } from '@/api/websocket/websocket'
import { States } from '@/models/Game'
import { useGameStore } from '@/stores/useGameStore'

type MessageType = 'ENGINE_ACTION'

type Action =
  | 'GetElement'
  | 'Continue'
  | 'TradeOffer'
  | 'RemoveTradeOffer'
  | 'TradeRequest'
  | 'TradeAck'

export abstract class BaseStateHandler<T extends States['State']> {
  protected readonly requiredState: T

  constructor (protected ws: WebsocketConnector, requiredState: T) {
    this.requiredState = requiredState
  }

  getState (): string {
    return this.requiredState
  }

  isValid (): boolean {
    const gameStore = useGameStore()
    if (gameStore.gameState.State !== this.requiredState) {
      return false
    }
    return true
  }

  protected ensureCorrectState (): void {
    const gameStore = useGameStore()
    if (gameStore.gameState.State !== this.requiredState) {
      throw new Error(
        `Операция недопустима. Требуемое состояние: ${this.requiredState}, текущее: ${gameStore.gameState.State}`
      )
    }
  }

  // Шаблонный метод для отправки команды с проверкой состояния
  protected sendCommand (
    content: { Type: MessageType; Action: Action } & Record<string, any>
  ): void {
    this.ensureCorrectState()
    this.ws.Send({
      // Type: 'ENGINE_ACTION',
      ...content
    })
  }
}
