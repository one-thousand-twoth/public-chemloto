<script setup lang="ts">
import { Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
const emit = defineEmits<{
    selectPlayer: [name: string]
}>()

</script>
<template>
    <div>
        <h2>Топ игроков</h2>
        <ul class="list-none p-0 font-bold m-0">
            <li @click="emit('selectPlayer', pl.Name)" class="break-words flex justify-between items-center p-2
                    hover:underline 
                    [&:nth-child(1)]:bg-amber-300
                    [&:nth-child(2)]:bg-stone-300
                    [&:nth-child(3)]:bg-yellow-500
                    rounded-md my-2"
                v-for="pl in GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score) "
                :key="pl.Name">
                {{ pl.Name }} - {{ pl.Score }}
            </li>
        </ul>
    </div>
</template>