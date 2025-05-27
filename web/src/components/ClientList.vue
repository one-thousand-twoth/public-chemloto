<script setup lang="ts">
import IconButton from '@/components/UI/IconButton.vue';
import IconButtonBackground from '@/components/UI/IconButtonBackground.vue';
import { i18nRole, Role, User } from '@/models/User';
import { useUsersListStore } from '@/stores/useUserListStore';
import { useUserStore } from '@/stores/useUserStore';
import { ArrowPathIcon, TrashIcon } from "@heroicons/vue/24/outline";
import { storeToRefs } from 'pinia';

const userStore = useUserStore()
const selfuser = userStore.getUser()
const usersListStore = useUsersListStore()
const { UsersList } = storeToRefs(usersListStore)
usersListStore.fetchUsers()
function Patch(usr: User) {
  let role = ""
  if (usr.role == Role.Player) {
    role = Role.Judge
  } else if (usr.role == Role.Judge) {
    role = Role.Player
  }
  if (confirm(`Вы действительно хотите изменить роль ${usr.username} на ${role} ?`)) {
    userStore.PatchUser(usr)
    usersListStore.fetchUsers()
  }
}
function Delete(usr: string) {
  userStore.Remove(usr)
  usersListStore.fetchUsers()

}
</script>
<template>
  <div class="p-4 bars shadow-lg bg-white">
    <table class="mb-4">
      <thead>
        <tr>
          <th>Имя</th>
          <th>Роль</th>
          <th>Комната</th>
          <th v-if="selfuser.role != Role.Player"></th>
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
          <td class="">{{ user.username }} <span v-if="selfuser.username == user.username">(Это
              вы!)</span></td>
          <td class="">{{ i18nRole(user.role) }}</td>
          <td> {{ user.room ? user.room : "-" }}</td>
          <td v-if="selfuser.role == Role.Admin">
            <div class="flex justify-end items-end gap-1">
              <button v-if="user.role == Role.Player" @click="Patch(user)">Назначить Судьей</button>
              <button v-if="user.role == Role.Judge" @click="Patch(user)">Cделать Игроком</button>
              <IconButtonBackground class="bg-red-700" :icon="TrashIcon" @click="Delete(user.username)">
              </IconButtonBackground>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="mb-5">
      <div class=" flex flex-row gap-2">
        <IconButton :icon="ArrowPathIcon" @click="usersListStore.fetchUsers()" />
      </div>
    </div>
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