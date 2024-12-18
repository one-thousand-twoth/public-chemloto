<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { Role } from '@/models/User';
import { Player, useGameStore } from '@/stores/useGameStore';
import { storeToRefs } from 'pinia';
import { computed, inject, ref } from 'vue';
// const props = defineProps<{
// 	players: Array<Player>;

// }>()
const GameStore = useGameStore()
const { gameState } = storeToRefs(GameStore)
const players = computed(() => gameState.value.Players.filter((k, _) => k.Role === Role.Player))
console.log(players)
const ws = inject('connector') as WebsocketConnector
interface TradeStruct {
	Player1: Player,
	Player2: Player,
	Element1: string,
	Element2: string,
}
const struct = ref({
	Player1: {
		Name: '',
		Role: Role.Player,
		Score: 0,
		RaisedHand: false,
		Bag: {},
		CompletedFields: []
	},
	Player2: {
		Name: '',
		Role: Role.Player,
		Score: 0,
		RaisedHand: false,
		Bag: {},
		CompletedFields: []
	},
	Element1: '',
	Element2: '',
})
function Trade(st: TradeStruct) {
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "TradeAdmin",
		Player1: st.Player1.Name,
		Player2: st.Player2.Name,
		Element1: st.Element1,
		Element2: st.Element2,

	})
}

</script>
<template>
	<form @submit.prevent="Trade(struct)">
		<div class="flex flex-wrap">
			<section class="flex min-w-96 flex-col gap-4 border-r-2 text-right border-gray-700  p-3">
				<label for="roomName">Игрок:</label>
				<select v-model="struct.Player1">
					<option disabled value="">Выберите</option>
					<option v-for="[_, pl] in Object.entries(players)" :value="pl">{{ pl.Name }}</option>
				</select>
				<label for="roomName">Элемент:</label>
				<select v-model="struct.Element1">
					<option disabled value="" class="text-right">Выберите</option>
					<option class="text-right"
						v-for="[i, _] in Object.entries(struct.Player1.Bag).filter(([k, _]) => k !== 'TRADE')">{{ i
						}}
					</option>
				</select>
			</section>
			<section class="flex flex-col gap-4 p-3">
				<label>Игрок:</label>
				<select v-model="struct.Player2">
					<option disabled value="">Выберите</option>
					<option v-for="[_, pl] in Object.entries(players)" :value="pl">{{ pl.Name }}</option>
				</select>
				<label>Элемент:</label>
				<select v-model="struct.Element2">
					<option disabled value="">Выберите</option>
					<option v-for="[i, _] in Object.entries(struct.Player2.Bag).filter(([k, _]) => k !== 'TRADE')">
						{{ i }}</option>
				</select>
			</section>
		</div>
		<button type="submit">Отправить</button>
	</form>
</template>

<style scoped>
select {
	width: 20rem;
}
</style>
