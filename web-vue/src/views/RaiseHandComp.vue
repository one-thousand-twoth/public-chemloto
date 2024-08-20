<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import ChemicalElementFormInput from '@/components/UI/ChemicalElementFormInput.vue';
import { Player } from '@/stores/useGameStore';
import { computed, inject, ref } from 'vue';
import polymers from '@/../../polymers.json'

const props = defineProps<{
	player: Player;
	
}>()

const ws = inject('connector') as WebsocketConnector

interface CheckStruct {
	Field: string,
	Name: string,
	// Structure: { [id: string]: number; }
}

function Check(ch: CheckStruct, str: { [id: string]: number; }) {
	console.log(str)
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "RaiseHand",
		Field: ch.Field,
		Name: ch.Name,
		Structure: Object.fromEntries(Object.entries(str).filter(([_, v]) => v !== 0)),
	})
	
}
const check = ref<CheckStruct>({
	Field: 'Альфа',
	Name: ''
	// Structure: {"C":34,"C6H4":14,"C6H5":16,"CH":21,"CH2":21,"CH3":26,"Cl":14,"H":46,"N":16,"O":23,"TRADE":3},
})

const struct = ref<{ [id: string]: number; }>(
	Object.fromEntries(Object.entries(props.player.Bag).map(([name, value]) => { return [name, 0] }))
)
console.log("str", struct.value)

interface Field {
	[key: string]: {
		[polymerName: string]: Polymer;
	};
}
interface Polymer extends Array<Entry> {
	// [PolymerName: string]: Array<any>;
}
interface Entry {
	[element: string]: number;
}

const polymersObj = Object.entries(polymers as Field)
const poly = polymers as Field
console.log(polymersObj)
console.log(poly)

</script>

<template>
	<form @submit.prevent="Check(check, struct)">
		<div class="flex flex-col gap-4 ">
			<section>
				<label for="roomName">Поле:</label>
				<select v-model="check.Field">
					<option disabled value="">Выберите</option>
					<option v-for="[field, k] in Object.entries(poly)">{{ field }}</option>
				</select>
			</section>
			<section v-if='check.Field'>
				<label for="roomName">Название структуры:</label>
				<select v-model="check.Name">
					<option disabled value="">Выберите</option>
					<option v-for="[v, k] in Object.entries(poly[check.Field])">{{ v }}</option>
				</select>
			</section>
			<div v-if="poly[check.Field][check.Name] !== undefined" class="flex flex-wrap justify-between">
				<ChemicalElementFormInput v-for="[elname, _] in Object.entries(poly[check.Field][check.Name][0])"
					:elname="elname" :max="player.Bag[elname]" v-model.number="struct[elname]" />
			</div>
			<button type="submit">Отправить</button>
		</div>
	</form>

</template>
