<script setup lang="ts">
import { IconRole, Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
const emit = defineEmits<{
    selectPlayer: [name: string]
}>()

</script>
<template>
    <div>
        <ul class="list-none p-0 font-bold m-0">
            <li @click="emit('selectPlayer', pl.Name)" class="block break-words text-xs flex flex-wrap justify-between items-center p-1 px-2 border-2 border-b-4 border-playing
                    hover:underline 
                    rounded-large my-2"
                v-for="pl in GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score)"
                :key="pl.Name">
                <div class=" inline-flex gap-1 items-center" v-bind="$attrs">
                    <component :is="IconRole(pl.Role)" class="size-4 text-slate-700" />
                    <span class="text-slate-700"> {{ pl.Name }}</span>
                </div>

                <span>{{ pl.Score }}</span>
            </li>
        </ul>
    </div>
</template>