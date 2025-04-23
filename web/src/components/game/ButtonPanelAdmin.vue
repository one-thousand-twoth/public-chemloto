<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { DesignButton } from '@/components/game';
import { Modal } from '@/components/UI/index';
import { ObtainStateHandler, TradeStateHandler, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ArrowDownCircleIcon,
    HandRaisedIcon, PuzzlePieceIcon, ShoppingBagIcon
} from "@heroicons/vue/24/solid";
import { storeToRefs } from "pinia";
import { computed, inject, ref } from 'vue';
import IconButton from '../UI/IconButton.vue';
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


const selectedTool = ref('puzzle')
const selectedPalace = ref<"strip" | "list">('strip')

function swap() {
    if (selectedPalace.value == 'strip') {
        selectedPalace.value = 'list'
    } else {
        selectedPalace.value = 'strip'
    }
}

</script>
<template>
    <div class="flex gap-2">

        <div class="flex  border-slate-300 hover:bg-slate-100 border-b-main border-b-2  shadow-large items-center cursor-pointer px-2 py-2 rounded border bg-white
              border-main-dark text-main
              border-slate-300 "
            @click="swap()">
            <component :is="selectedPalace == 'strip' ? ShoppingBagIcon : ArrowDownCircleIcon" class="size-7 lg:size-10 text-slate-500" />
        </div>


        <div class="flex rounded shadow-large border border-b-main border-b-2">
            <DesignButton class="rounded-none rounded-l" v-model="selectedTool" value="puzzle" label="Puzzle">
                <PuzzlePieceIcon class="size-7 lg:size-10" />
            </DesignButton>

            <DesignButton class="rounded-none rounded-r" v-model="selectedTool" value="hand" label="Hand">
                <HandRaisedIcon class="size-7 lg:size-10 -rotate-90" />
            </DesignButton>
        </div>
    </div>

    <!-- <template v-if="gameStore.gameState.Status === 'STATUS_WAITING'">
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
</Modal> -->
</template>