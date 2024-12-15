<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { hasTimer, useGameStore } from '@/stores/useGameStore';
import { PauseIcon, PlayIcon, StopIcon } from '@heroicons/vue/24/solid';
import { computed, inject } from 'vue';
import IconButton from './IconButton.vue';
const GameStore = useGameStore();

const timerString = computed(() => {
    if (!GameStore.timer || GameStore.timer == null) {
        return `0 : 0`
    }
    if (GameStore.timer < 0) {
        return `0 : 0`
    }
    let minutes = Math.floor(GameStore.timer / 60);
    let remainingSeconds = GameStore.timer % 60;

    return `${minutes} : ${remainingSeconds}`
})

const ws = inject('connector') as WebsocketConnector

function TimerStop() {
    ws.Send(
        {
            "Type": "ENGINE_ACTION",
            "Action": "TimerStop",
        }
    )
}

function TimerPlay() {
    ws.Send(
        {
            "Type": "ENGINE_ACTION",
            "Action": "TimerPlay",
        }
    )
}

function TimerPause() {
    ws.Send(
        {
            "Type": "ENGINE_ACTION",
            "Action": "TimerPause",
        }
    )
}

</script>

<template>
    <div class="flex relative items-center">
        <div class="px-4 py-4 whitespace-nowrap text-3xl font-3xl text-gray-800">
            {{ timerString }}
        </div>
        <template v-if="hasTimer(GameStore.gameState)">
            <IconButton class="" :icon="StopIcon" @click="TimerStop"></IconButton>
            <IconButton v-if="GameStore.gameState.StateStruct.TimerStatus == 'Stopped'" class="" :icon="PlayIcon"
                @click="TimerPlay"></IconButton>
            <IconButton v-if="GameStore.gameState.StateStruct.TimerStatus == 'Started'" class="" :icon="PauseIcon"
                @click="TimerPause"></IconButton>
            <div>
                {{ GameStore.gameState.StateStruct.TimerStatus }}
            </div>
        </template>
    </div>  
</template>