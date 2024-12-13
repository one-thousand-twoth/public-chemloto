<script setup lang="ts">
import { RaiseHandComp, TradeExchange, UserElements } from '@/components/game/';
import { Modal } from '@/components/UI/index';
import { useGameStore } from '@/stores/useGameStore';
import { computed, ref } from 'vue';

// const props = defineProps<{
//     modal: string;
// }>()

// const ws = inject('connector') as WebsocketConnector

const gameStore = useGameStore()
// const userStore = useUserStore()


const currPlayer = computed(() => {
    return gameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

const curInfoPlayer = ref('')
// const curCheckPlayer = ref<Hand>()
const RaiseHandButton = ref(false)
const TradeButton = ref(false)



</script>
<template>
    <template v-if="gameStore.gameState.Status === 'STATUS_WAITING'">
        <button disabled >Ждем начала</button>
    </template>
    <template v-else>
        <template v-if="gameStore.gameState.State == 'OBTAIN'">
            <button @click="RaiseHandButton = !RaiseHandButton">Поднять
                руку</button>
        </template>
        <template v-if="gameStore.gameState.State == 'HAND'">
            <button  @click="RaiseHandButton = !RaiseHandButton">Поднять
                руку</button>
        </template>
        <template v-if="gameStore.gameState.State == 'TRADE'">
            <button @click="TradeButton = !TradeButton">Обменять</button>
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
        <template v-if="gameStore.SelfPlayer" #body>
            <RaiseHandComp :player="gameStore.SelfPlayer" />
        </template>
    </Modal>
    <Modal :show="TradeButton" @close="TradeButton = false">
        <template #header>
            <h3 class="font-bold text-center">Обменять</h3>
        </template>
        <template v-if="gameStore.SelfPlayer" #body>
            <TradeExchange />
        </template>
    </Modal>
</template>