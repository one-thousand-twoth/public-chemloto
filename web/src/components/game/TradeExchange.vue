<script setup lang="ts">
import { ElementImage, IconButton } from '@/components/UI';
import { GameInfo, StateTRADE, TradeStateHandler, useGameStore } from '@/stores/useGameStore';
import {
    // ChatBubbleOvalLeftEllipsisIcon,
    CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";

import { storeToRefs } from 'pinia';
import { computed, ref } from 'vue';
import TradeExchangeStocks from './TradeExchangeStocks.vue';


const isTradingState = computed(() => gameState.value.State === "TRADE");


const gameStore = useGameStore()
const { gameState } = storeToRefs(gameStore)
const player = gameStore.SelfPlayer

interface TradeOffer {
    Element: string
    toElement: string
}

const tradeForm = ref({
    Element: "",
    toElement: "",
})

const tradeState = computed(() => {
    if (gameState.value.State === "TRADE") {
        return gameState.value as GameInfo & StateTRADE
    }
    return null
})

const tradeHandler = computed(() => {
    if (!isTradingState.value) return null;
    return gameStore.currentStateHandler as TradeStateHandler;
});

function handleTrade(offer: TradeOffer) {
    if (!tradeHandler.value) return;
    tradeHandler.value.trade(offer.Element, offer.toElement);
}
function cancelTrade() {
    if (!tradeHandler.value) return;
    tradeHandler.value.cancelTrade();
}

function ackTrade(req: string) {
    if (!tradeHandler.value) return;
    tradeHandler.value.ackTrade(req);
}

const selfStock = computed(() => tradeState.value?.StateStruct?.StockExchange.StockList.find(stock => stock.Owner === player.Name) ?? null)

const stockList = computed(() => {
    if (!tradeState.value?.StateStruct?.StockExchange.StockList) return [];
    return Object.entries(tradeState.value.StateStruct.StockExchange.StockList).filter(([_, v]) => v.Owner !== player.Name);
});
const requests = computed(() => {
    if (!selfStock.value) return []
    return Object.entries(selfStock.value?.Requests).filter(([_, v]) => v.Accept)
})
const alreadyTradedStruct = computed(() => {
    return tradeState.value?.StateStruct?.StockExchange.TradeLog.find(log => log.User === player.Name) ?? null
})

</script>
<template>
    <div class="w-full md:w-[60vw]">
        <div v-if="!isTradingState" class="w-full text-lg text-gray-600">
            Биржа недоступна
        </div>

        <div v-else class="w-full flex flex-wrap justify-center gap-4">
            <div v-if="alreadyTradedStruct" class="px-4 py-2 flex items-center justify-center border-solid border-2 border-blue-400 rounded-lg;">
                На этом ходу вы поменялись:
                <div class=" text-lg inline-flex gap-1 items-center">
                    <ElementImage class=" w-8 inline m-1" :elname="alreadyTradedStruct.GetElement" />
                    <span>за</span>
                    <ElementImage class="w-8 inline m-1" :elname="alreadyTradedStruct.GaveElement" />
                </div>
            </div>
            <form v-else-if="!selfStock" class="flex  flex-col gap-4" @submit.prevent="handleTrade(tradeForm)">

                <div>
                    <section class="flex flex-col gap-1 mb-2 items-end">
                        <label>Элемент:</label>
                        <select v-model="tradeForm.Element">
                            <option disabled value="" class="text-right">Выберите</option>
                            <option class="text-right"
                                v-for="[i, _] in Object.entries(player.Bag).filter(([k, _]) => k !== 'TRADE')">{{ i
                                }}
                            </option>
                        </select>
                        <label>Обменять на:</label>
                        <select v-model="tradeForm.toElement">
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
            </form>

            <div class="w-[max(20rem)]" v-else>
                <div class="flex flex-nowrap flex-col mb-1">
                    <span class="mb-1">Ваш лот:</span>
                    <div>
                        <div class=" w-[min(20rem)] border-solid border-2 border-blue-400 rounded-lg px-4 py-2">
                            <div class="flex">
                                <div class=" text-lg inline-flex gap-1 items-center">
                                    <ElementImage class=" w-8 inline m-1" :elname="selfStock.Element" />
                                    <span>за</span>
                                    <ElementImage class="w-8 inline m-1" :elname="selfStock.ToElement" />
                                </div>
                                <IconButton @click="cancelTrade" class="ml-auto" :icon="XMarkIcon" />
                            </div>

                        </div>
                    </div>
                </div>
                <details>
                    <summary>
                        <div
                            class="inline-flex w-5 h-5 text-[1rem] items-center justify-center rounded-full bg-blue-400 text-white  font-bold">
                            {{ Object.entries(selfStock.Requests).length }}
                        </div>
                        Согласны:
                    </summary>
                    <!--                <IconButton class="" :icon="XMarkIcon" />-->
                    <p class="inline-flex items-center w-full" v-for="([_, request]) in requests">
                        {{ request.Player }}
                        <IconButton class="ml-auto" @click="ackTrade(request.ID)" :icon="CheckIcon" />
                    </p>
                </details>
            </div>
            <div>
                <TradeExchangeStocks :stockList="stockList" :tradeHandler="tradeHandler" />
            </div>
        </div>
    </div>
</template>

<style scoped>
select {
    width: 20rem;
}
</style>
