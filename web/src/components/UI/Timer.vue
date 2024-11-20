<script setup lang="ts">
import { useGameStore } from '@/stores/useGameStore';
import { computed } from 'vue';

const GameStore = useGameStore();

const timerString = computed(() => {
    if (!GameStore.timer || GameStore.timer == null){
        return `0 : 0`
    }
    if (GameStore.timer < 0) {
         return `0 : 0`
    }
    let minutes = Math.floor(GameStore.timer / 60);
    let remainingSeconds = GameStore.timer % 60;
    if (GameStore.timer > 1) {
        setTimeout(() => {
            if (GameStore.timer ==null){
                return
            }
            GameStore.timer--; // Decrement the timer count
        }, 1000);
    }
    return `${minutes} : ${remainingSeconds}`
})

</script>

<template>
    <div class="px-6 py-4 whitespace-nowrap text-3xl font-3xl text-gray-800">
        {{ timerString }}
    </div>
</template>