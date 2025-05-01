<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ElementImage, IconButtonChecked, UserInfo } from '@/components/UI';
import { StockEntity } from '@/models/Game';
import { Role } from '@/models/User';
import TradeStateController from '@/state_controllers/trade';
import { useGameStore } from '@/stores/useGameStore';
import {
	// ChatBubbleOvalLeftEllipsisIcon,
	CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";
import { ChatBubbleOvalLeftEllipsisIcon } from "@heroicons/vue/24/solid";
import { computed, inject } from 'vue';


const ws = inject('connector') as WebsocketConnector
const gameStore = useGameStore()
// const { gameState } = storeToRefs(gameStore)
const player = gameStore.SelfPlayer!

interface TradeRequest {
	StockID: string
	Accept: boolean
}
defineProps<{
	stockList: [string, StockEntity][];
}>()

const tradeController = new TradeStateController(ws)

function requestTrade(req: TradeRequest) {
	if (!tradeController.isValid()) return;
	tradeController.requestTrade(req.StockID, req.Accept);
}
function accepted(stock: StockEntity) {
	return stock.Requests[player.Name] && stock.Requests[player.Name].Accept
}
function notAccepted(stock: StockEntity) {
	return stock.Requests[player.Name] && !stock.Requests[player.Name].Accept
}
</script>

<template>
	<div class="">
		<div v-if="stockList.length === 0">
			Другие игроки ещё не сделали своих предложений...
		</div>
		<div class="flex text-sm flex-nowrap mb-2 flex-col" v-for="[_, Stock] in stockList">
			<div class="mb-1 flex gap-1">
				<UserInfo :name="Stock.Owner" :role="Role.Player" />
				<span> предлагает:</span>
			</div>
			<div class=" w-fit  border-solid border-2 border-blue-400 rounded-lg px-2 py-1">
				<div class="flex">
					<div class=" text-lg inline-flex gap-1 items-center">
						<ElementImage class=" w-8 inline m-1" :elname="Stock.Element" />
						<span>за</span>
						<ElementImage class="w-8 inline m-1" :elname="Stock.ToElement" />
					</div>
					<div class="ml-auto"></div>

					<ChatBubbleOvalLeftEllipsisIcon v-if="accepted(Stock)"
						class="relative left-2 bottom-1 transform scale-x-[-1] h-6 w-6 text-blue-400" />
					<IconButtonChecked :is-checked="accepted(Stock)" @click="requestTrade({
						StockID: Stock.ID,
						Accept: true,
					})" class=" " :icon="CheckIcon" />
					<IconButtonChecked :is-checked="false"
						:class="{ 'border-2 rounded-lg border-red-500': notAccepted(Stock) }" @click="requestTrade({
							StockID: Stock.ID,
							Accept: false,
						})" :icon="XMarkIcon" />
				</div>
			</div>
		</div>
	</div>
</template>
