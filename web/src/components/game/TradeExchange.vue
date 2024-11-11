<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ElementImage, IconButton, UserInfo } from '@/components/UI/';
import { Role } from '@/models/User';
import { GameInfo, Player, StateTRADE, StockEntity, TradeStateHandler, useGameStore } from '@/stores/useGameStore';
import { useUserStore } from '@/stores/useUserStore';
import {
    ChatBubbleOvalLeftEllipsisIcon,
    CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";
import { storeToRefs } from 'pinia';
import { computed, inject, onMounted, ref } from 'vue';


const isTradingState = computed(() => gameState.value.State === "TRADE");

const gameStore = useGameStore()
const { gameState } = storeToRefs(gameStore)
const player = gameStore.SelfPlayer

interface TradeOffer {
    Element: string
    toElement: string
}
interface TradeRequest {
    StockID: string
    Accept: boolean
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
function requestTrade(req: TradeRequest) {
    if (!tradeHandler.value) return;
    tradeHandler.value.requestTrade(req.StockID, req.Accept);
}
function accepted(stock: StockEntity){
    // return stock.Requests
    console.log(stock.Requests);
    console.log(Array.isArray(stock.Requests))
    return Object.entries(stock.Requests).find(([name, _]) => name === player.Name ) ?? null
}
// tradeState.value?.Trade()
const selfStock = computed(() => tradeState.value?.StateStruct?.StockExchange.StockList.find(stock => stock.Owner === player.Name) ?? null)

const stockList = computed(() => {
    if (!tradeState.value?.StateStruct?.StockExchange.StockList) return [];
    return Object.entries(tradeState.value.StateStruct.StockExchange.StockList);
});

</script>
<template>
    <div v-if="!isTradingState" class="text-lg text-gray-600">
        Биржа недоступна
    </div>

    <div v-else class="flex flex-wrap gap-4">
        <form v-if="!selfStock" class="flex  flex-col gap-4  p-3" @submit.prevent="handleTrade(tradeForm)">

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
                        {{ 1 }}
                    </div>
                    Согласны:
                </summary>
                <IconButton class="" :icon="XMarkIcon" />
                <p>
                    Requires a computer running an operating system. The computer must have some
                    memory and ideally some kind of long-term storage. An input device as well
                    as some form of output device is recommended.Requires a computer running an operating system. The
                    computer must have some
                    memory and ideally some kind of long-term storage. An input device as well
                    as some form of output device is recommended.Requires a computer running an operating system. The
                    computer must have some
                    memory and ideally some kind of long-term storage. An input device as well
                    as some form of output device is recommended.Requires a computer running an operating system. The
                    computer must have some
                    memory and ideally some kind of long-term storage. An input device as well
                    as some form of output device is recommended.Requires a computer running an operating system. The
                    computer must have some
                    memory and ideally some kind of long-term storage. An input device as well
                    as some form of output device is recommended.
                </p>
            </details>
        </div>
        <div>
            <div class="flex flex-nowrap mb-2   flex-col" v-for="[_, Stock] in stockList">
                <div class="mb-1 flex gap-1">
                    <UserInfo :name="Stock.Owner" :role="Role.Player" />
                    <span> предлагает:</span>
                </div>
                <div class=" w-[min(20rem)] border-solid border-2 border-blue-400 rounded-lg px-4 py-2">

                    <div class="flex">
                        <div class=" text-lg inline-flex gap-1 items-center">
                            <ElementImage class=" w-8 inline m-1" :elname="Stock.Element" />
                            <span>за</span>
                            <ElementImage class="w-8 inline m-1" :elname="Stock.ToElement" />
                        </div>
                        <div class="ml-auto" ></div>
                        
                        <ChatBubbleOvalLeftEllipsisIcon v-if="accepted(Stock)" class="relative left-2 bottom-1 transform scale-x-[-1] h-6 w-6 text-gray-500" />
                        <IconButton @click="requestTrade({
                            StockID: Stock.ID,
                            Accept: true,
                        })" :icon="CheckIcon" />
                        <IconButton class="" :icon="XMarkIcon" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
select {
    width: 20rem;
}
</style>
