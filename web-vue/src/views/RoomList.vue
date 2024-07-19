<script lang="ts">
import { defineComponent, computed, ref } from "vue"
import { useRoomsStore } from '../stores/useRoomsStore'
import { IconButton } from '../components/UI/index'
import {
  ArrowPathIcon,
  LinkIcon,
  PaintBrushIcon,
  UserIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
} from "@heroicons/vue/24/outline";
// import { useToasterStore } from "../stores/useToasterStore";
// import store from "./store/index"
// const toasterStore = useToasterStore();
export default defineComponent({
  components: { IconButton },
  setup() {
    const roomStore = useRoomsStore()
    const rooms = computed(() => roomStore.roomList)
    const showModal = ref(false)
    const showCreateModal = ref(false)

    const createGameInput = ref('')

    // const onDeleteNote = (noteId: number) => {
    //   showModal.value = true
    //   store.commit("setCurrentId", noteId)
    // }
    // const confirmDelete = (permission: boolean): void => {
    //   if (permission) {
    //     store.commit("deleteNote")
    //   }
    //   closeModal()
    // }
    // function closeModal() {
    //   store.commit("setCurrentId", 0)
    //   showModal.value = false
    // }

    return {
      rooms,
      roomStore,
      ArrowPathIcon,
      showModal,
      createGameInput,
      // closeModal,
      // confirmDelete
    }
  }
})
</script>
<template>
  <div style="position: relative;">
    <div class="mb-5 ">
      <table >
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
              <button>Подключиться</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="mb-5">
      <div class=" flex flex-row gap-2">
        <input type="text" v-model="createGameInput" class="" />
        <button @click="roomStore.CreateGame(createGameInput)">Создать</button>
        <IconButton :icon="ArrowPathIcon" @click="roomStore.Fetch()" />
      </div>
    </div>
  </div>
</template>



<style scoped>
table {
  border-collapse: collapse;
  width: 100%;
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