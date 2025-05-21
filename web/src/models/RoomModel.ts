export type Room = {
  name: string
  type: string
  engine: { Status: EngineStatus }
}

export type CreateRoomRequest = {
  name: string
} & EngineConfig

type EngineConfig = PolymersConfig | AminoConfig

type PolymersConfig = {
  type: 'polymers'
  engineConfig: {
    maxPlayers: number
    elementCounts: { [el: string]: number }
    time: number
    isAuto: boolean
    isAutoCheck: boolean
  }
}
type AminoConfig = {
  type: 'amino'
  engineConfig: {}
}

export type EngineStatus =
  | 'STATUS_WAITING'
  | 'STATUS_STARTED'
  | 'STATUS_COMPLETED'

export function i18nStatus (status: EngineStatus): string {
  switch (status) {
    case 'STATUS_WAITING':
      return 'Ожидает'
    case 'STATUS_STARTED':
      return 'Запущена'
    case 'STATUS_COMPLETED':
      return 'Завершена'
    default: {
      const _exhaustiveCheck: never = status
      return _exhaustiveCheck
    }
  }
}
