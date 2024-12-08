<script lang="ts">
import { WebsocketConnector } from "@/api/websocket/websocket";
import CreateRoom from "@/components/CreateRoom.vue";
import Modal from "@/components/UI/Modal.vue";
import { Role } from "@/models/User";
import { useUserStore } from "@/stores/useUserStore";
import {
  ArrowPathIcon
} from "@heroicons/vue/24/outline";
import { computed, defineComponent, inject, ref } from "vue";
import { useRoomsStore } from '../stores/useRoomsStore';
import { IconButton } from './UI/index';
export default defineComponent({
  components: { IconButton, Modal, CreateRoom },
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
      ConnectGame,
    }
  }
})
</script>
<template>
  <!-- <div class="h-screen flex items-center justify-center bg-opacity-5 bg-slate-500">
    <div class=" flex items-center  justify-center gap-4 flex-col"> -->
  <div class="p-4 shadow-lg bg-white">
    <table class="mb-4">
      <thead>
        <tr>
          <th>Имя</th>
          <th>Статус</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="roomStore.fetching">
          <td colspan="3">Загрузка...</td>
        </tr>
        <tr v-else-if="!rooms.length">
          <td colspan="3">Пока нет доступных комнат</td>
        </tr>
        <tr v-else v-for="room in rooms" :key="room.name">
          <td class="">{{ room.name }}</td>
          <td></td>
          <td>
            <button @click="ConnectGame(room.name)">Подключиться</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="mb-5">
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