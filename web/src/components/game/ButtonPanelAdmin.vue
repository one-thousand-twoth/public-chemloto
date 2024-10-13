<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import RaiseHandComp from '@/components/game/RaiseHandComp.vue';
import { Modal } from '@/components/UI/index';
import { Hand, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import { computed, inject, ref } from 'vue';
import Trade from './Trade.vue';

// const props = defineProps<{
//     modal: string;
// }>()

const ws = inject('connector') as WebsocketConnector

const GameStore = useGameStore()
const userStore = useUserStore()


function StartGame() {
    ws.Send({
        Type: 'HUB_STARTGAME',
        Name: userStore.UserCreds!.room.toString()
    })
}
function GetElement() {
    console.log('Get element!')
    ws.Send({
        Type: 'ENGINE_ACTION',
        Action: 'GetElement'
    })
}
function SendContinue() {
    ws.Send({
        Type: 'ENGINE_ACTION',
        Action: 'Continue'
    })
}

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
const curCheckPlayer = ref<Hand>()
const RaiseHandButton = ref(false)
const TradeButton = ref(false)



</script>
<template>
    <template v-if="!GameStore.gameState.Started">
        <button @click="StartGame()">
            Начать игру
        </button>
    </template>
    <template v-else>
        <template v-if="GameStore.gameState.State == 'OBTAIN'">
            <button @click="GetElement()">Достать
                элемент</button>
        </template>
        <template v-if="GameStore.gameState.State == 'HAND'">
            <button disabled @click="GetElement()">Ждем
                проверки</button>
        </template>
        <template v-if="GameStore.gameState.State == 'TRADE'">
            <button @click="TradeButton = !TradeButton">Обменять</button>
            <button @click="SendContinue()">Продолжить</button>
        </template>
    </template>
    <Modal :show="TradeButton" @close="TradeButton = false">
        <template #header>
            <h3 class="font-bold text-center">Обменять</h3>
        </template>
        <template v-if="GameStore.SelfPlayer" #body>
            <Trade :players="GameStore.gameState.Players" />
        </template>
    </Modal>
</template>