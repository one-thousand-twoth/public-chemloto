<script lang="ts">
import { defineComponent, computed, ref } from "vue"
import { useRoomsStore } from '../store/index'
import { IconButton} from '../components/UI/index'
import {
  ArrowPathIcon,
  LinkIcon,
  PaintBrushIcon,
  UserIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
} from "@heroicons/vue/24/outline";
// import store from "./store/index"
export default defineComponent({
  components: { IconButton },
  setup() {
    const counterStore = useRoomsStore()
    const rooms = computed(() => counterStore.RoomList)

    const showModal = ref(false)

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
      ArrowPathIcon,
      // onDeleteNote,
      showModal,
      // closeModal,
      // confirmDelete
    }
  }
})
</script>
<template>

    <h3 class="font-bold text-center">Подключиться к игре</h3>

    <table>
      <thead>
        <tr>
          <th>Имя</th>
          <th>Статус</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="room in rooms" :key="room.name">
          <!-- <td>{{ ID }}</td>
          <td>{{ name }}</td>
          <td>{{ IP }}</td>
          <td>{{ status }}</td> -->
          <td>{{ room.name }}</td>
          <td>{{ room.status }}</td>
          <td>
            <button>Подключиться</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="flex flex-row gap-2">
      <button>Создать</button>
      <IconButton :icon="ArrowPathIcon" />
    </div>
</template>



<style>
.completed {
  text-decoration: line-through;
}
</style>