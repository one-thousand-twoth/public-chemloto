<script setup lang="ts">
import { useGameStore } from '@/stores/useGameStore';
import { computed, watch } from 'vue';

const GameStore = useGameStore();

const timerString = computed(() => {
    var minutes = Math.floor(GameStore.timer / 60);
    var remainingSeconds = GameStore.timer % 60;
    if (GameStore.timer > 1) {
        setTimeout(() => {
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