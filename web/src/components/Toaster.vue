<template>
    <Teleport to="body">
        <div v-if="toastStore.toasts.length"
            class="absolute flex flex-col gap-4 right-4 bottom-4 rounded-lg shadow-lg bg-white z-50">
            <div v-for="(toast, idx) in toastStore.toasts" :key="idx" class="p-4 flex flex-row gap-2">
                <ExclamationCircleIcon v-if="toast.status === 'error'" class="size-6 text-red-500" />
                <InformationCircleIcon v-if="toast.status === 'info'" class="size-6 text-blue-500" />
                {{ toast.text }}
            </div>
        </div>
    </Teleport>
</template>
<script setup lang="ts">
import obtain from "@/assets/sounds/obtain.wav";
import {
    ExclamationCircleIcon,
    InformationCircleIcon,
} from "@heroicons/vue/24/outline";

import { useToasterStore } from "../stores/useToasterStore";
const toastStore = useToasterStore();
let audio = new Audio(obtain);
toastStore.callback = () => {
    audio.play();
}
</script>