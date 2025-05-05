<script lang="ts">
import { WebsocketConnector } from "@/api/websocket/websocket";
import CreateRoom from "@/components/CreateRoom.vue";
import { ElementImage, IconButton, Modal } from '@/components/UI';
import { i18nStatus } from "@/models/RoomModel";
import { Role } from "@/models/User";
import { useUserStore } from "@/stores/useUserStore";
import {
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";
import {
  PlayIcon
} from "@heroicons/vue/24/solid";
import { computed, defineComponent, inject, ref } from "vue";
import { useRoomsStore } from '../stores/useRoomsStore';
export default defineComponent({
  components: { IconButton, Modal, CreateRoom, ElementImage, PlayIcon },
  setup() {
    const userStore = useUserStore()
    const selfuser = userStore.getUser()
    const roomStore = useRoomsStore()
    const rooms = computed(() => roomStore.roomList)
    const showModal = ref(false)
    const ws = inject('connector') as WebsocketConnector
    roomStore.Fetch()
    function ConnectGame(roomName: string) {
      ws.Send(
        {
          "Type": "HUB_SUBSCRIBE",
          "Target": "room",
          "Name": roomName
        }
      )
    }
    return {
      rooms,
      Role,
      selfuser,
      roomStore,
      userStore,
      ArrowPathIcon,
      showModal,
      i18nStatus,
      ConnectGame,
    }
  }
})
</script>
<template>
  <!-- <div class="h-screen flex items-center justify-center bg-opacity-5 bg-slate-500">
    <div class=" flex items-center  justify-center gap-4 flex-col"> -->
  <div class="p-4 shadow-lg  max-w-[80%] mx-auto ">

    <div class="flex flex-col gap-1 justify-stretch ">
      <div v-for="room in rooms" :key="room.name"  class="flex gap-2 bg-white shadow-large bars px-2 py-1">
        <div class="flex gap-2 flex-wrap w-full" >
          <ElementImage class="w-10" elname="H" />
          <div>
            <p class="text-sm">Полимеры</p> 
            <p class="text-lg"> {{ room.name }}</p>
          </div>
          <div class="mx-auto my-auto w-40">
            {{ i18nStatus(room.engine.Status) }}
          </div>
        </div>
        <button class="ml-auto my-1" @click="ConnectGame(room.name)"><PlayIcon class="size-6"/></button>
      </div>
    </div>
   
    <div class="mt-4">
      <div class=" flex flex-row gap-2">
        <button v-if="selfuser.role != Role.Player" @click="showModal = !showModal">Создать</button>
        <IconButton :icon="ArrowPathIcon" @click="roomStore.Fetch()" />
      </div>
    </div>

    <Modal :show="showModal" @close="showModal = !showModal">
      <template #header>
        <h3 class="font-bold text-center">Создать игру</h3>
      </template>
      <template #body>
        <CreateRoom></CreateRoom>
      </template>
    </Modal>
  </div>
  <!-- </div> -->

  <!-- </div> -->
</template>



<style scoped>
table {
  border-collapse: collapse;
}

table {
  display: flex;
  flex-flow: column;
  width: 100%;
  height: 400px;

}

thead {
  padding-right: 13px;
  flex: 0 0 auto;
}

tbody {
  flex: 1 1 auto;
  display: block;
  overflow-y: scroll;
  overflow-x: scroll;
}

tr {
  width: 100%;
  display: table;
  table-layout: fixed;
}
</style>