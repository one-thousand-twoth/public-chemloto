<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ElementImage, IconButtonChecked, UserInfo } from '@/components/UI';
import { StockEntity } from '@/models/Game';
import { Role } from '@/models/User';
import TradeStateController from '@/state_controllers/trade';
import { useInterfaceStore } from '@/stores/RoomInterface';
import { useGameStore } from '@/stores/useGameStore';
import {
	// ChatBubbleOvalLeftEllipsisIcon,
	CheckIcon, XMarkIcon
} from "@heroicons/vue/24/outline";
import { ChatBubbleOvalLeftEllipsisIcon } from "@heroicons/vue/24/solid";
import { storeToRefs } from 'pinia';
import { computed, inject } from 'vue';


const ws = inject('connector') as WebsocketConnector
const gameStore = useGameStore()

const InterfaceStore = useInterfaceStore()
const { currentPlayerSelection } = storeToRefs(InterfaceStore)

// const { gameState } = storeToRefs(gameStore)
const player = gameStore.SelfPlayer!

interface TradeRequest {
	StockID: string
	Accept: boolean
}
const props = defineProps<{
	stockList: [string, StockEntity][];
}>()

const SelectedStock = computed(()=>{return props.stockList.find(([_, v])=>{return v.Owner == currentPlayerSelection.value})?.[1]})

const tradeController = new TradeStateController(ws)

function requestTrade(req: TradeRequest) {
	if (!tradeController.isValid()) return;
	tradeController.requestTrade(req.StockID, req.Accept);
}
function accepted(stock: StockEntity) {
	return (stock.Requests[player.Name] && stock.Requests[player.Name].Accept) ?? false
}
function notAccepted(stock: StockEntity) {
	return stock.Requests[player.Name] && !stock.Requests[player.Name].Accept
}
</script>

<template>
	<div class="">
		<div class="bg-gray-200 py-2 rounded-[2px] text-base text-center flex-col flex items-center justify-center h-full md:text-lg text-gray-600"
		 v-if="SelectedStock === undefined">
			Выберите игрока в панели слева, чтобы отобразить его предложение
		</div>
        <div v-else>
            <div class="inline-flex  flex-wrap items-center justify-end mb-1 flex gap-1">
					<UserInfo :name="SelectedStock.Owner" :role="Role.Player" />
					<span> отдает:</span>
				</div>
				<div class="flex">
					<div class=" text-lg inline-flex gap-1 items-center">
						<ElementImage class=" w-8 inline m-1" :elname="SelectedStock.Element" />
						<span>за</span>
						<ElementImage class="w-8 inline m-1" :elname="SelectedStock.ToElement" />
					</div>
					<div class="ml-auto"></div>

					<ChatBubbleOvalLeftEllipsisIcon v-if="accepted(SelectedStock)"
						class="relative left-2 bottom-1 transform scale-x-[-1] h-6 w-6 text-blue-400" />
					<IconButtonChecked :is-checked="accepted(SelectedStock)" @click="requestTrade({
						StockID: SelectedStock.ID,
						Accept: true,
					})" class=" " :icon="CheckIcon" />
					<IconButtonChecked :is-checked="false"
						:class="{ 'border-2 rounded-lg border-red-500': notAccepted(SelectedStock) }" @click="requestTrade({
							StockID: SelectedStock.ID,
							Accept: false,
						})" :icon="XMarkIcon" />
				</div>
        </div>
	</div>
</template>
