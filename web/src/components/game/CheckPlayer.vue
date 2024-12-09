<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import ChemicalElementFormInput from '@/components/UI/ChemicalElementFormInput.vue';
import { Hand } from '@/stores/useGameStore';
import { inject, ref } from 'vue';
import { Polymers } from '@/models/Polymers';
const props = defineProps<{
	player: Hand;
}>()

const ws = inject('connector') as WebsocketConnector

interface CheckStruct {
	Field: string,
	Name: string,
}
function Check(Player: string, accept: boolean) {
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "Check",
		Player: Player,
		Accept: true
	})
}
const check = ref<CheckStruct>({
	Field: props.player.Field,
	Name: props.player.Name
})
const struct = ref<{ [id: string]: number; }>(
	Object.fromEntries(Object.entries(props.player.Structure))
)
console.log("str", struct.value)


</script>

<template>
	<form @submit.prevent>
		<div class="flex flex-col gap-4 ">
			<div class="text-lg">Поле: {{ player.Field }}</div>
			<div class="text-lg">Структура: {{ player.Name }}</div>
			<div v-if="Polymers[check.Field][check.Name] !== undefined" class="flex flex-wrap justify-between">
				<ChemicalElementFormInput :disabled="true"
					v-for="[elname] in Object.entries(Polymers[check.Field][check.Name][0])" :elname="elname"
					:max="player.Player.Bag[elname]" v-model.number="struct[elname]" />
			</div>
			<div class="  inline-flex gap-1 justify-between">
				<button @click="Check(player.Player.Name, true)">Принять</button>
				<button class="bg-red-500" @click="Check(player.Player.Name, false)">Отклонить</button>
			</div>
		</div>
	</form>
</template>
