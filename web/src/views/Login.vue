<script setup lang="ts">
import { useUserStore } from '@/stores/useUserStore';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const username = ref("")
const code = ref("")
const formErrors = ref<Record<string, string>>({})

const userStore = useUserStore();

async function onSubmit() {
    formErrors.value = {}
    const result = await userStore.Login(username.value, code.value);

    if (result?.formErrors) {
        // Устанавливаем ошибки формы для отображения
        formErrors.value = result.formErrors;
        return;
    }

    if (result?.success && userStore.UserCreds) {
        console.log(await router.replace({ name: "Hub" }))
    }
}

const router = useRouter()

const isChecked = ref(false)

function Check() {
    isChecked.value = !isChecked.value
    code.value = ""
}

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
                        <label>Введите название команды:
                            <input v-model.trim="username" class="w-full" autocomplete="off" type="text" maxlength="25" required>
                        </label>
                        <p v-if="formErrors.name" class="text-red-500 text-sm mt-1">{{ formErrors.name }}</p>
                        <label>
                            <input type="checkbox" v-on:click="Check()" />
                            Я админ
                        </label>
                        <label v-show="isChecked">Введите код:
                            <input :required="isChecked ? true : undefined" v-model="code" autocomplete="off"
                                type="text" maxlength="25">
                        </label>
                        <p v-if="formErrors.code" class="text-red-500 text-sm mt-1">{{ formErrors.code }}</p>
                        <button type="submit">Войти</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>