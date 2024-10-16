<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ElementImage } from '@/components/UI/';
import { Role } from '@/models/User';
import { Player, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import { storeToRefs } from 'pinia';
import { computed, inject, onMounted, ref } from 'vue';

const GameStore = useGameStore()
const { gameState } = storeToRefs(GameStore)
const player = GameStore.SelfPlayer
const ws = inject('connector') as WebsocketConnector
interface TradeOffer {
    Element: string
    toElement: string
}
const struct = ref({
    Element: "",
    toElement: "",
})
function Trade(st: TradeOffer) {
    ws.Send({
        Type: "ENGINE_ACTION",
        Action: "TradeOffer",
        Element: st.Element,
        ToElement: st.toElement
    })
}


const checkGameState = () => {
    if (gameState.value.State !== 'TRADE') {
        throw new Error(`Invalid state: Expected 'TRADE', but got '${gameState.value.State}'`);
    }
};

onMounted(() => {
    checkGameState();
});

</script>
<template>
    <div class="grid grid-cols-2 gap-4">


        <form class="flex flex-col gap-4 border-r-2 text-right border-gray-700  p-3" @submit.prevent="Trade(struct)">
            <div >
                <section class="flex flex-col items-end">
                    <label for="roomName">Элемент:</label>
                    <select v-model="struct.Element">
                        <option disabled value="" class="text-right">Выберите</option>
                        <option class="text-right"
                            v-for="[i, _] in Object.entries(player.Bag).filter(([k, _]) => k !== 'TRADE')">{{ i
                            }}
                        </option>
                    </select>
                    <label>Обменять на:</label>
                    <select v-model="struct.toElement">
                        <option disabled value="" class="text-right">Выберите</option>
                        <option class="text-right"
                            v-for="[i, _] in Object.entries(gameState.Bag.Elements).filter(([k, _]) => k !== 'TRADE')">
                            {{ i
                            }}
                        </option>
                    </select>
                </section>

            </div>
            <button type="submit" class="self-end">Отправить</button>

        </form>
        <div v-if="gameState.State === 'TRADE'">
            <div v-for="[_, Stock] in Object.entries(gameState.StateStruct!.StockExchange.StockList)">
                <span>{{ Stock.Owner }} предлагает</span>
                <ElementImage  class="w-8 inline m-1" :elname="Stock.Element" />
                <span>за</span>
                <ElementImage class="w-8 inline m-1" :elname="Stock.ToElement" />

            </div>
        </div>
    </div>
</template>

<style scoped>
select {
    width: 20rem;
}
</style>
