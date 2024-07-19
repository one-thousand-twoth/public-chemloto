<script setup lang="ts">
import { useUserStore } from '@/stores/useUserStore';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const username = ref("")

const userStore = useUserStore();

async function onSubmit(){
    await userStore.createUser(username.value);
    if (userStore.UserCreds) {
       console.log( await router.push({name: "RoomList"}))
    }
}

const router = useRouter()

</script>
<template>
    <div class="h-screen flex items-center justify-center ">
        <div class=" flex items-center justify-center gap-4 flex-col">
            <div class="p-4 shadow-lg ">
                <div>
                    <h2>Добро пожаловать на турнир по "Химлото"</h2>
                </div>
                <div>
                    <form @submit.prevent="onSubmit()" class="w-full flex items-start justify-center flex-col gap-4">
                        <h2>Вход</h2>
                        <label for="name">Введите имя:</label>
                        <input v-model="username" autocomplete="off" type="text" maxlength="25" required style="width: 95%;">
                        <button type="submit">Войти</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>