<script setup lang="ts">
import { Role, UserInfo, useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from 'pinia';
import { ArrowPathIcon, TrashIcon } from "@heroicons/vue/24/outline";
import IconButton from './UI/IconButton.vue';
import IconButtonBackground from './UI/IconButtonBackground.vue';

const userStore = useUserStore()
const { UsersList } = storeToRefs(userStore)
userStore.fetchUsers()
function Patch(usr: UserInfo) {
  let role = ""
  if (usr.role == Role.Player) {
    role = Role.Judge
  } else if (usr.role == Role.Judge) {
    role = Role.Player
  }
  if (confirm(`Вы действительно хотите изменить роль ${usr.username} на ${role} ?`)){
    userStore.PatchUser(usr)
    userStore.fetchUsers()
  }
}
function Delete(usr: string){
  userStore.Remove(usr)
  userStore.fetchUsers()

}
</script>
<template>
  <div class="h-screen flex items-center justify-center bg-opacity-5 bg-slate-500">
    <div class=" flex items-center  justify-center gap-4 flex-col">
      <div class="p-4 shadow-lg bg-white" style="width: 70%;">
        <table class="mb-4">
          <thead>
            <tr>
              <th>Имя</th>
              <th>Роль</th>
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
              <td class="">{{ user.role }}</td>
              <td> {{ user.room ? user.room : "-" }}</td>
              <td>
                <div class="flex justify-end items-end gap-1">
                  <button v-if="user.role == Role.Player" @click="Patch(user)">Назначить Судьей</button>
                  <button v-if="user.role == Role.Judge" @click="Patch(user)">Cделать Игроком</button>
                  <IconButtonBackground class="bg-red-700" :icon="TrashIcon" @click="Delete(user.username)"></IconButtonBackground>
                </div>
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