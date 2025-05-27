<script setup lang="ts">
import { PlayIcon, PlusCircleIcon, XCircleIcon } from '@heroicons/vue/24/solid'
import { useElementSize } from '@vueuse/core'
import { computed, inject, nextTick, onMounted, ref, useTemplateRef } from 'vue'

import { WebsocketConnector } from '@/api/websocket/websocket'
import CreateRoom from '@/components/CreateRoom.vue'
import { ElementImage } from '@/components/UI'
import { i18nStatus } from '@/models/RoomModel'
import { Role } from '@/models/User'
import { useUserStore } from '@/stores/useUserStore'
import { useRoomsStore } from '../stores/useRoomsStore'

// Store setup
const userStore = useUserStore()
const selfuser = userStore.getUser()
const roomStore = useRoomsStore()
const rooms = computed(() => roomStore.roomList)

// Component state
const showModal = ref(false)
const createRoomRef = useTemplateRef("createRoom")
// @ts-ignore
const { height: refHeight } = useElementSize(createRoomRef)
const isAuto = ref(false)

const collapsedStyle = computed(() => {
  if (showModal.value) {
    return `height: ${refHeight.value}px`
  }
  return ''
})

// Websocket
const ws = inject('connector') as WebsocketConnector

// Fetch rooms on component mount
roomStore.Fetch()

// Functions
function ConnectGame(roomName: string) {
  ws.Send({
    "Type": "HUB_SUBSCRIBE",
    "Target": "room",
    "Name": roomName
  })
}
async function toggleCreateRoom() {
  if (showModal.value === true) {
    isAuto.value = false

  }
  showModal.value = !showModal.value

}
function transitionend() {
  if (showModal.value === true) {
    isAuto.value = true
  }
}
onMounted(() => {
  const el = createRoomRef.value
  if (el) {
    const height = el.$el.offsetHeight
    console.log('Element height:', height)
  }
})
</script>

<template>
  <div class="">
    <div class="flex flex-1 flex-col gap-2 justify-stretch">
      
      <!-- Окно создания комнаты -->
      <div v-if="selfuser?.hasPermission()"
        class="flex flex-col bars border-dotted outline outline-main outline-2 px-2 py-2">
        <div @click="toggleCreateRoom()">
          <div class="flex gap-2 border-dotted px-2 py-1" :class="!showModal ? '' : 'bg-slate-200'">
            <div :class="!showModal ? '' : 'opacity-0'" class="m-auto transition-all">Создать комнату</div>
            <PlusCircleIcon v-show="!showModal" class="w-10 flex-shrink-0 text-gray-500" />
            <XCircleIcon v-show="showModal" class="w-10 flex-shrink-0 text-gray-500" />
          </div>
        </div>
        <div @transitionend="transitionend()" :style="collapsedStyle"
          class="w-full h-0 transition-all duration-300 ease-in overflow-hidden">
          <CreateRoom ref="createRoom" @exit-panel="toggleCreateRoom()" />
        </div>
      </div>
      
      <!-- Если комнат нет, placeholder -->
      <div v-if="rooms.length === 0" class="flex items-center justify-center">
        Пока комнат нет
      </div>
      
      <!-- Блоки комнат -->
      <div v-for="room in rooms" :key="room.name" class="flex gap-2 bg-white shadow-large bars px-2 py-1">
        <div class="flex gap-2 items-center flex-wrap w-full">
          <div class="flex gap-2 flex-1">
            <ElementImage class="w-10 shrink-0" elname="H" />
            <div class="w-28 flex-1">
              <p class="text-sm">Полимеры</p>
              <p class="text-lg  flex-grow truncate whitespace-nowrap overflow-hidden text-ellipsis">{{ room.name }}</p>
            </div>
          </div>
          <div class="w-20">
            {{ i18nStatus(room.engine.Status) }}
          </div>
        </div>
        <button class="ml-auto my-1" @click="ConnectGame(room.name)">
          <PlayIcon class="size-6" />
        </button>
      </div>
      <div class="h-20"></div>
    </div>
  </div>
</template>

<style scoped>
.collapsed {
  height: 0px;
}
</style>