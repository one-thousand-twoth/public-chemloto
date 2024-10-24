<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ElementImage, IconButton } from '@/components/UI/';
import { Role } from '@/models/User';
import { GameInfo, Player, StateTRADE, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";
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

const tradeState = computed(()=>{
    if (gameState.value.State ==="TRADE"){
        return gameState.value as GameInfo & StateTRADE
    }
    return null
})

const selfStock = computed(() => tradeState.value?.StateStruct?.StockExchange.StockList.find( stock => stock.Owner === player.Name) ?? null)

</script>
<template>
    <div v-if="gameState.State === 'TRADE'" class="grid grid-cols-2 gap-4">


        <form class="flex flex-col gap-4 border-r-2 text-right border-gray-700  p-3" @submit.prevent="Trade(struct)">

            <div v-if="!selfStock">
                <section class="flex flex-col gap-1 mb-2 items-end">
                    <label>Элемент:</label>
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
                <button type="submit" class="self-end">Отправить</button>
            </div>
            <div v-else>

            </div>

        </form>
        <div v-if="gameState.State === 'TRADE'">
            <div class="flex flex-nowrap" v-for="[_, Stock] in Object.entries(gameState.StateStruct!.StockExchange.StockList)">
                <div>
                    <div>
                        <span>{{ Stock.Owner }} предлагает:</span>
                    </div>
                    <ElementImage class="w-8 inline m-1" :elname="Stock.Element" />
                    <span>за</span>
                    <ElementImage class="w-8 inline m-1" :elname="Stock.ToElement" />
                </div>
                <IconButton class=" ml-auto" :icon="CheckIcon" />
                <IconButton class="" :icon="XMarkIcon" />
            </div>
        </div>
    </div>
</template>

<style scoped>
select {
    width: 20rem;
}
</style>
