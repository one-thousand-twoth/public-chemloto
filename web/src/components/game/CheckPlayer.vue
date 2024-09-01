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
function Check(ch: CheckStruct, str: { [id: string]: number; }, Player: string) {
	console.log(str)
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "Check",
		Player: Player,
		Field: ch.Field,
		Name: ch.Name,
		Structure: Object.fromEntries(Object.entries(str).filter(([_, v]) => v !== 0)),
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
	<form @submit.prevent="Check(check, struct, player.Player.Name)">
		<div class="flex flex-col gap-4 ">
			<section>
				<label for="roomName">Поле:</label>
				<select v-model="check.Field">
					<option disabled value="">Выберите</option>
					<option v-for="[field] in Object.entries(Polymers)">{{ field }}</option>
				</select>
			</section>
			<section v-if='check.Field'>
				<label for="roomName">Название структуры:</label>
				<select v-model="check.Name">
					<option disabled value="">Выберите</option>
					<option v-for="[v, _] in Object.entries(Polymers[check.Field])">{{ v }}</option>
				</select>
			</section>
			<div v-if="Polymers[check.Field][check.Name] !== undefined" class="flex flex-wrap justify-between">
				<ChemicalElementFormInput v-for="[elname] in Object.entries(Polymers[check.Field][check.Name][0])"
					:elname="elname" :max="player.Player.Bag[elname]" v-model.number="struct[elname]" />
			</div>
			<button type="submit">Отправить</button>
		</div>
	</form>
</template>
