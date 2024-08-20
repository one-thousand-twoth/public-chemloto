<script setup lang="ts">
import { useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from 'pinia';
import { ArrowPathIcon } from "@heroicons/vue/24/outline";
import IconButton from './UI/IconButton.vue';

const userStore = useUserStore()
const { UsersList } = storeToRefs(userStore)
userStore.fetchUsers()
</script>
<template>
  <div class="h-screen flex items-center justify-center bg-opacity-5 bg-slate-500">
    <div class=" flex items-center  justify-center gap-4 flex-col">
      <div class="p-4 shadow-lg bg-white" style="width: 70%;">
        <table class="mb-4">
          <thead>
            <tr>
              <th>Имя</th>
              <th>Комната</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="userStore.fetching">
              <td colspan="3">Загрузка...</td>
            </tr>
            <tr v-else-if="!UsersList.length">
              <td colspan="3">Пока нет пользователей</td>
            </tr>
            <tr v-else v-for="user in UsersList" :key="user.username">
              <td class="">{{ user.username }}</td>
              <td> {{ user.room ? user.room : "-"}}</td>
              <td>
                <button>Редактировать</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div class="mb-5">
          <div class=" flex flex-row gap-2">
            <!-- <button @click="showModal = !showModal">Создать</button> -->
            <IconButton :icon="ArrowPathIcon" @click="userStore.fetchUsers()" />
          </div>
        </div>
      </div>
    </div>
    <!-- <Modal :show="showModal" @close="showModal = !showModal">
        <template #header>
          <h3 class="font-bold text-center">Создать игру</h3>
        </template>
<template #body>
          <CreateRoom></CreateRoom>
        </template>
</Modal> -->
  </div>
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