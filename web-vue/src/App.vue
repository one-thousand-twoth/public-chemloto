<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { Client } from "./api/core/client";
import RoomList from './views/RoomList.vue'
import { Navbar, Modal ,IconButton} from './components/UI/index'
import {
  ArrowPathIcon,
  LinkIcon,
  PaintBrushIcon,
  UserIcon,
  ExclamationCircleIcon,
  InformationCircleIcon,
} from "@heroicons/vue/24/outline";



const gamesModalVisible = ref(false);
const UserModalVisible = ref(false);
const authPanelVisible = ref(false);
const token = ref(localStorage.getItem("token") ?? "");

if (!token.value) {
  UserModalVisible.value = true;
}

type ToastStatus = "error" | "info" | "ok";

function toast(message: string, type: ToastStatus) {
  toasts.value.push({ message, type });
}
const toasts = ref<{ message: string; type: ToastStatus }[]>([]);

const client = new Client("http://localhost:1090/", token.value);

const username = ref(localStorage.getItem("username") ?? "");

const currentGameID = ref("-");

// const timer = ref<number | null>(null);
// onMounted(() => {
//   timer.value = setInterval(() => {
//     toasts.value.shift();
//   }, 3000);
// });

// onUnmounted(() => {
//   if (timer.value) clearInterval(timer.value);
// });


async function createUser() {
  const resp = await fetch(client.url(`/api/v1/users`), {
    method: "POST",
    // headers: client.headers(),
    body: new URLSearchParams({
      name: username.value,
      // IP: "172.16.1.126",
    }),
  });

  if (resp.status === 409) {
    toast(`Пользователь с именем ${username.value} уже существует!`, "error");
    return;
  }

  if (!resp.ok) {
    console.error("Failed to login with user");
    toast(`Не удалось войти под именем ${username.value}`, "error");
    return;
  }

  token.value = await resp.text();
  localStorage.setItem("token", token.value);
  localStorage.setItem("username", username.value);
  client.token = token.value;

  toast(`Вы вошли под пользователем ${username.value}`, "info");
  console.log(`Token for user ${username.value} created: ${token.value}`);
  UserModalVisible.value = false;
}
async function logout() {
  toast("Вы вышли из учетной записи", "info");
  token.value = "";
  username.value = "";
  localStorage.clear();
  UserModalVisible.value = true;
}
</script>



<template>
  <Navbar>
    <div class="m-auto font-bold flex gap-2">
        <div class="m-auto">{{ currentGameID }}</div>
        <IconButton
          @click="gamesModalVisible = !gamesModalVisible"
          :icon="LinkIcon"
        />
        <IconButton
        @click="authPanelVisible = !authPanelVisible"
        :icon="UserIcon"
      />
    <div v-if="token">
      <div class="font-bold">{{ username }}</div>
      <button @click="logout()">Выйти</button>
    </div>
    <!-- <div v-else class="flex flex-row gap-2">
      <input type="text" v-model="username" />
      <button>Войти</button>
    </div> -->
  </div>
  </Navbar>
  <Modal :show="gamesModalVisible">

    <template #header>
      <h3>Custom Header</h3>
    </template>

    <template #body>
      <RoomList class="min-h-96" />
    </template>

  </Modal>
  <Modal :show="UserModalVisible">

    <template #header>
      <h3 class="font-bold">Введите своё имя</h3>
    </template>

    <template #body>
      <div class="flex flex-row gap-2">
        <input type="text" v-model="username" />
        <button @click="createUser()">Войти</button>
      </div>
    </template>

    <template #footer>
      <span></span>
    </template>

  </Modal>
  <div class="absolute flex flex-col gap-4 right-4 bottom-4 rounded-lg shadow-lg bg-white z-10">
    <div v-for="(toast, idx) in toasts" :key="idx" class="p-4 flex flex-row gap-2">
      <ExclamationCircleIcon v-if="toast.type === 'error'" class="size-6 text-red-500" />
      <InformationCircleIcon v-if="toast.type === 'info'" class="size-6 text-blue-500" />
      {{ toast.message }}
    </div>
  </div>
</template>

<style scoped>
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
</style>
