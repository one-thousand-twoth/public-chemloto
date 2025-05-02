<script setup lang="ts">
import { ElementImage, IconButton } from '@/components/UI';
import { useGameStore } from '@/stores/useGameStore';
import {
    // ChatBubbleOvalLeftEllipsisIcon,
    CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";

import { WebsocketConnector } from '@/api/websocket/websocket';
import { GameInfo, StateTRADE } from '@/models/Game';
import { TradeStateHandler } from '@/state_controllers';
import { storeToRefs } from 'pinia';
import { computed, inject, ref } from 'vue';
import TradeExchangeStocks from './TradeExchangeStocks.vue';
import TradeSelected from './TradeSelected.vue';

const gameStore = useGameStore()
const { gameState, SelfPlayer } = storeToRefs(gameStore)
const player = SelfPlayer

const ws = inject('connector') as WebsocketConnector

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

const tradeController = new TradeStateHandler(ws)

function handleTrade(offer: TradeOffer) {
    if (!tradeController.isValid) return;
    tradeController.trade(offer.Element, offer.toElement);
}
function cancelTrade() {
    if (!tradeController.isValid) return;
    tradeController.cancelTrade();
}

function ackTrade(req: string) {
    if (!tradeController.isValid) return;
    tradeController.ackTrade(req);
}

const selfStock = computed(() => tradeState.value?.StateStruct?.StockExchange.StockList.find(stock => stock.Owner === player.value?.Name) ?? null)

const stockList = computed(() => {
    if (!tradeState.value?.StateStruct?.StockExchange.StockList) return [];
    return Object.entries(tradeState.value.StateStruct.StockExchange.StockList).filter(([_, v]) => v.Owner !== player.value?.Name);
});
const requests = computed(() => {
    if (!selfStock.value) return []
    return Object.entries(selfStock.value?.Requests).filter(([_, v]) => v.Accept)
})
const alreadyTradedStruct = computed(() => {
    return tradeState.value?.StateStruct?.StockExchange.TradeLog.find(log => log.User === player.value?.Name) ?? null
})

console.log("TradeState", tradeController.isValid())

</script>
<template>
    <div class="w-full h-full overflow-y-auto ">

        <div v-if="!tradeController.isValid()"
            class="w-full bg-gray-200 rounded-[2px] text-center flex-col flex items-center justify-center h-full text-lg text-gray-600">
            Торговля ещё не началась
            <ElementImage class=" w-8 inline m-1" elname="TRADE" />
            <p class="loading"></p>
        </div>

        <div v-else class="relative h-full w-full flex flex-col  gap-4">
            <div v-if="alreadyTradedStruct"
                class="px-4 py-2 flex items-center justify-center border-solid border-2 border-blue-400 rounded-lg;">
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
                                v-for="[i, _] in Object.entries(player!.Bag).filter(([k, _]) => k !== 'TRADE')">{{ i
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

            <div class=" flex flex-col h-1/2 text-sm  bars pb-2 border-0 border-b-2 " v-else>
                <div class="flex flex-nowrap flex-col mb-1">
                    <!-- <span class="mb-1">Ваш лот:</span> -->
                    <div>
                        <div class="w-full border-solid border-2 border-blue-400 rounded-lg px-2 py-1">
                            <div class="flex ">
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
                <div class="relative flex flex-col px-2 grow overflow-y-auto">
                    <div>
                        <div
                            class="inline-flex w-5 h-5 text-sm items-center justify-center rounded-full bg-blue-400 text-white font-medium">
                            {{Object.entries(selfStock.Requests).filter(([_, v]) => { return v.Accept == true }).length}}
                        </div>
                        Согласны:
                    </div>
                    <div class="">
                        <p class=" inline-flex items-center w-full" v-for="([_, request]) in requests">
                            {{ request.Player }}
                            <IconButton class="ml-auto" @click="ackTrade(request.ID)" :icon="CheckIcon" />
                        </p>

                    </div>
                    <div v-if="requests.length == 0" class="block my-auto text-center ">
                        Еще никто из игроков не согласился на ваше предложение
                    </div>
                </div>

            </div>
            <div>
                <TradeSelected :stockList="stockList" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.loading:after {
    overflow: hidden;
    display: inline-block;
    vertical-align: bottom;
    -webkit-animation: ellipsis steps(4, end) 900ms infinite;
    animation: ellipsis steps(4, end) 900ms infinite;
    content: "\2026";
    /* ascii code for the ellipsis character */
    width: 0px;
}

@keyframes ellipsis {
    to {
        width: 1.25em;
    }
}

@-webkit-keyframes ellipsis {
    to {
        width: 1.25em;
    }
}
</style>
