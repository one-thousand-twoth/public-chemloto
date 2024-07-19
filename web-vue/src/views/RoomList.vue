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
    const rooms = computed(() => roomStore.RoomList)
    const showModal = ref(false)
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
  <div  style="position: relative;">
    <table class="table-fixed">
      <thead>
        <tr>
          <th class="min-w-30">Имя</th>
          <th>Статус</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="roomStore.Fetching">
          <td colspan="3">Загрузка...</td>
        </tr>
        <tr v-else-if="!rooms.length">
          <td colspan="3">Пока нет доступных комнат</td>
        </tr>
        <tr v-else v-for="room in rooms" :key="room.name">
          <!-- <td>{{ ID }}</td>
          <td>{{ name }}</td>
          <td>{{ IP }}</td>
          <td>{{ status }}</td> -->
          <td class="min-w-80">{{ room.name }}</td>
          <!-- <td>{{ room.status }}</td> -->
          <td>
            <button>Подключиться</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="absolute bottom-4">
    <input type="text" v-model="createGameInput" class="mb-2"
     
      />
    <div class=" flex flex-row gap-2">
      <button>Создать</button>
      <IconButton :icon="ArrowPathIcon" @click="roomStore.Fetch()" />
    </div>
  </div>
  </div>
</template>



<style>
.completed {
  text-decoration: line-through;
}
</style>