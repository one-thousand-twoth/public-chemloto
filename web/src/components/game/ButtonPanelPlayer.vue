<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import RaiseHandComp from '@/components/game/RaiseHandComp.vue';
import { Modal } from '@/components/UI/index';
import { Hand, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import { computed, inject, ref } from 'vue';

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
        <button disabled >Ждем начала</button>
    </template>
    <template v-else>
        <template v-if="GameStore.gameState.State == 'OBTAIN'">
            <button @click="RaiseHandButton = !RaiseHandButton">Поднять
                руку</button>
        </template>
        <template v-if="GameStore.gameState.State == 'HAND'">
            <button  @click="RaiseHandButton = !RaiseHandButton">Поднять
                руку</button>
        </template>
        <template v-if="GameStore.gameState.State == 'TRADE'">
        </template>
    </template>

    <Modal :show="currPlayer !== undefined" @close="curInfoPlayer = ''">
        <template #header>
            <h3 class="font-bold text-center">Информация о игроке {{ curInfoPlayer }}</h3>
        </template>
        <template #body>
            <UserElements v-if="currPlayer" :player="currPlayer" />
        </template>
    </Modal>
    <Modal :show="RaiseHandButton" @close="RaiseHandButton = false">
        <template #header>
            <h3 class="font-bold text-center">Поднять руку</h3>
        </template>
        <template v-if="GameStore.SelfPlayer" #body>
            <RaiseHandComp :player="GameStore.SelfPlayer" />
        </template>
    </Modal>
</template>