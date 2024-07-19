<script setup lang="ts">
import { ref } from "vue";
import { Client } from "./api/core/client";
import RoomList from './views/RoomList.vue'
import { Navbar, Modal, IconButton, Tooltip } from './components/UI/index'
import Toaster from './components/Toaster.vue'
import { useToasterStore } from "./stores/useToasterStore";
import { useUserStore } from "./stores/useUserStore";
import {
  ArrowPathIcon,
  LinkIcon,
  PaintBrushIcon,
  UserIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
} from "@heroicons/vue/24/outline";
import { storeToRefs } from "pinia";
import UserLogin from './components/UserLogin.vue'



const gamesModalVisible = ref(false);
const UserModalVisible = ref(false);
const authPanelVisible = ref(false);
const exitPanelVisible = ref(false);
// const token = ref(localStorage.getItem("token") ?? "");


// const client = new Client("http://localhost:1090/", token.value);

const currentGameID = ref("-");

const toasterStore = useToasterStore();
const userStore = useUserStore();
const { UserCreds } = storeToRefs(userStore)


if (!userStore.UserCreds) {
  UserModalVisible.value = true;
}


async function logout() {
  exitPanelVisible.value = false;
  toasterStore.info("Вы вышли из учетной записи");
  userStore.$reset()
  UserModalVisible.value = true;
}


</script>



<template>
  <UserLogin />
  <div class="flex flex-col h-screen w-screen bg-[#fafafa]">
    <Navbar class="flex flex-row justify-between">
      <div class="m-auto font-bold flex gap-2">
        <div class="m-auto">{{ currentGameID }}</div>
        <IconButton @click="gamesModalVisible = !gamesModalVisible" :icon="LinkIcon" />
      </div>
      <IconButton @click="exitPanelVisible = !exitPanelVisible" :icon="UserIcon" />
      <Tooltip v-if="exitPanelVisible" class="right-4 top-12 left-auto">
        <div v-if="userStore.UserCreds">
          <div class="font-bold">{{ userStore.UserCreds.username }}</div>
          <button @click="logout()">Выйти</button>
        </div>
      </Tooltip>
    </Navbar>
  </div>
  <Modal :show="gamesModalVisible">

    <template #header>
      <h3 class="font-bold text-center">Подключиться к игре</h3>
    </template>

    <template #body>
      <RoomList />
    </template>

  </Modal>

  <Toaster />


</template>

<!-- <style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style> -->
