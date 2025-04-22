<script setup lang="ts">
import { UserInfo } from '@/components/UI/';
import { Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
const emit = defineEmits<{
    selectPlayer: [name: string]
}>()

</script>
<template>
    <div>
        <ul class="list-none p-0 font-bold m-0">
            <li @click="emit('selectPlayer', pl.Name)" class="break-words flex flex-wrap justify-between items-center p-2
                    hover:underline 
                    [&:nth-child(1)]:bg-amber-300
                    [&:nth-child(2)]:bg-stone-300
                    [&:nth-child(3)]:bg-yellow-500
                    rounded-md my-2"
                v-for="pl in GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score) "
                :key="pl.Name">
                <UserInfo :role="pl.Role" :name="pl.Name" />
                <span>{{ pl.Score }}</span>
            </li>
        </ul>
    </div>
</template>