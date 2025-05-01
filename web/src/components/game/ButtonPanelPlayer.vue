<script setup lang="ts">
import { DesignButton, RaiseHandComp, SelectedElements, TradeExchange, UserElements } from '@/components/game/';
import { Modal } from '@/components/UI';
import { useGameStore } from '@/stores/useGameStore';
import {
    ArrowDownCircleIcon,
    HandRaisedIcon, PuzzlePieceIcon, ShoppingBagIcon
} from "@heroicons/vue/24/solid";
import { openModal } from "jenesius-vue-modal";
import { computed, ref } from 'vue';

const selectedBtn = defineModel<"strip" | "list">('btn', { required: true })
const selectedRadio = defineModel<'puzzle' | 'trade'>('radio', { required: true })


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


// const selectedTool = ref('puzzle')
// const selectedBtn = ref('strip')

function swap() {
    if (selectedBtn.value == 'strip') {
        selectedBtn.value = 'list'
    } else {
        selectedBtn.value = 'strip'
    }
}

</script>
<template>
    <!-- <template v-if="gameStore.gameState.Status === 'STATUS_WAITING'">
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
</template> -->
    <div class="flex  gap-2">
        <!-- <button @click="openModal(NewModal, { msg: 'Welcome to jenesius-vue-modal' })">Hello</button> -->
        <div class="flex  border-slate-300 hover:bg-slate-100 border-b-main border-b-2  shadow-large items-center cursor-pointer px-2 py-2 rounded border bg-white
              border-main-dark text-main
              border-slate-300 " @click="swap">
            <component :is="selectedBtn == 'strip' ? ShoppingBagIcon : ArrowDownCircleIcon"
                class="size-7 lg:size-10 text-slate-500" />
        </div>

        <div class="flex rounded shadow-large border border-b-main border-b-2">
            <DesignButton class="rounded-none rounded-l" v-model="selectedRadio" value="puzzle" label="Puzzle">
                <PuzzlePieceIcon class="size-7 lg:size-10" />
            </DesignButton>

            <DesignButton class="rounded-none rounded-r" v-model="selectedRadio" value="trade" label="Trade">
                <HandRaisedIcon class="size-7 lg:size-10 -rotate-90" />
            </DesignButton>
        </div>
    </div>

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