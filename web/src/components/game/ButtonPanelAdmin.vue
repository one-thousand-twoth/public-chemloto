<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { Modal } from '@/components/UI/index';
import { ObtainStateHandler, TradeStateHandler, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from "pinia";
import { computed, inject, ref } from 'vue';
import Trade from './Trade.vue';

// const props = defineProps<{
//     modal: string;
// }>()

const ws = inject('connector') as WebsocketConnector

const gameStore = useGameStore()
const { gameState } = storeToRefs(gameStore)
const userStore = useUserStore()


function StartGame() {
    ws.Send({
        Type: 'HUB_STARTGAME',
        Name: userStore.UserInfo!.room
    })
}

const isObtainState = computed(() => gameState.value.State === "OBTAIN");
const isTradeState = computed(() => gameState.value.State === "TRADE");

const ObtainHandler = computed(() => {
    if (!isObtainState.value) return null;
    return gameStore.currentStateHandler as ObtainStateHandler;
});
const TradeHandler = computed(() => {
    if (!isTradeState.value) return null;
    return gameStore.currentStateHandler as TradeStateHandler;
});



function GetElement() {
    if (!ObtainHandler.value) return;
    ObtainHandler.value.getElement()
}
function sendContinue() {
    if (ObtainHandler.value) {
        ObtainHandler.value.sendContinue()
    }
    if (TradeHandler.value) {
        TradeHandler.value.sendContinue()
    }
    console.log("State doesn`t provide sendContinue action")
}

// const currPlayer = computed(() => {
//     return gameState.Players.find(player => player.Name === curInfoPlayer.value)
// })

// const curInfoPlayer = ref('')
// const curCheckPlayer = ref<Hand>()
// const RaiseHandButton = ref(false)
const TradeButton = ref(false)



</script>
<template>
    <template v-if="!gameStore.gameState.Started">
        <button @click="StartGame()">
            Начать игру
        </button>
    </template>
    <template v-else>
        <template v-if="gameStore.gameState.State == 'OBTAIN'">
            <button @click="GetElement()">Достать
                элемент</button>
        </template>
        <template v-if="gameStore.gameState.State == 'HAND'">
            <button disabled @click="GetElement()">Ждем
                проверки</button>
        </template>
        <template v-if="gameStore.gameState.State == 'TRADE'">
            <button @click="TradeButton = !TradeButton">Обменять</button>
            <button @click="sendContinue()">Продолжить</button>
        </template>
    </template>
    <Modal :show="TradeButton" @close="TradeButton = false">
        <template #header>
            <h3 class="font-bold text-center">Обменять</h3>
        </template>
        <template v-if="gameStore.SelfPlayer" #body>
            <Trade :players="gameStore.gameState.Players" />
        </template>
    </Modal>
</template>