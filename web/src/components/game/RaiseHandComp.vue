<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ChemicalElementCounter } from '@/components/game';
import { Player } from '@/models/Game';
import { Field, Polymer, Polymers, StructureNames } from '@/models/Polymers';
import { useGameStore } from '@/stores/useGameStore';
import { useKeyboardStore } from '@/stores/useRaiseHand';
import { storeToRefs } from 'pinia';
import { computed, inject, ref, watch } from 'vue';
import ChemicalElementCounterMobile from './ChemicalElementCounterMobile.vue';

const props = defineProps<{
	player: Player;
}>()

const gameStore = useGameStore()


const handStore = useKeyboardStore()
const { InputsValues } = storeToRefs(handStore)

const ws = inject('connector') as WebsocketConnector

interface CheckStruct<K extends Field> {
	Field: K,
	Name: StructureNames<K> | undefined
}

type InputEntries = { [id: string]: string; }

function Check(ch: CheckStruct<any>, str: InputEntries) {
	console.log(str)
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "RaiseHand",
		Field: ch.Field,
		Name: ch.Name,
		Structure: Object.fromEntries(Object.entries(str).map(([k, v]) => [k, Number(v)])),
	})
}

const isMobile = computed(() => {
	return navigator.maxTouchPoints > 1
})

const availableFields = Object.entries(Polymers)

console.debug(availableFields)

type AllCheckStructs = {
	[K in Field]: CheckStruct<K>
}[Field];

const check = ref<AllCheckStructs>({
	Field: 'Альфа',
	Name: undefined,
})


const struct = ref<InputEntries>(
	Object.fromEntries(Object.entries(props.player.Bag).map(([name]) => { return [name, '0'] }))
)

// const struct = Structure.value = Object.fromEntries(Object.entries(props.player.Bag).map(([name]) => { return [name, 0] }))


// type CurrenElements = typeof Polymers[check.Field][check.Name][0]

const selectedElem = defineModel<string | number>("selectedElem")

const currentElements = computed<Polymer | undefined>(() => {
	const { Field, Name } = check.value;
	if (!Name) return undefined;
	switch (Field) {
		case 'Альфа':
			return Polymers['Альфа'][Name];
		case 'Бета':
			return Polymers['Бета'][Name];
		case 'Гамма':
			return Polymers['Гамма'][Name];
		default:
			return undefined;
	}
});

watch(currentElements, () => {
	selectedElem.value = ''
	Object.keys(InputsValues.value).forEach((key) => {
		InputsValues.value[key] = '0';
	});


})


const isSubmitedHand = computed(() => {
	return gameStore.gameState.RaisedHands.some(
		(v) => v.Player.Name === gameStore.SelfPlayer?.Name
	)
})



</script>

<template>
	<form @submit.prevent="Check(check, struct)" class="flex h-full flex-col gap-2 mx-auto">
		<section class="flex ">
			<select v-model="check.Field" id="field"
				class=" w-auto bg-gray-50 border border-gray-300 text-gray-900 text-base rounded-l-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5>">
				<option value="Альфа" selected>α</option>
				<option value="Бета">β</option>
				<option value="Гамма">γ</option>
			</select>
			<select class="rounded-r-lg   w-full" v-model="check.Name">
				<option disabled value="">Выберите</option>
				<option v-for="[v] in Object.entries(Polymers[check.Field])">{{ v }}</option>
			</select>
		</section>
		<div v-if="currentElements !== undefined" class="flex  my-auto flex-wrap justify-around items-center">
			<ChemicalElementCounterMobile v-if="isMobile" :key="elname"
				v-for="elname in Object.keys(currentElements[0])" v-model:selected="selectedElem" :selector="elname"
				:elname="elname" :max="player.Bag[elname] ?? 0" v-model:input_value="struct[elname]" />
			<ChemicalElementCounter v-else :key="'_' + elname" v-for="elname in Object.keys(currentElements[0])"
				:elname="elname" :max="player.Bag[elname] ?? 0" v-model:input_value="struct[elname]">

			</ChemicalElementCounter>
		</div>
		<div v-else class="bg-main-tint h-full bars border-dashed flex  justify-center items-center text-sm"> Выберите структуру
			для сборки </div>
		<button v-if="!isSubmitedHand" class="self-end" type="submit">Отправить</button>
		<div class="mt-auto bg-main-tint px-2 py-1 rounded-[8px] loading " v-else> Дождитесь пока судья проверит
			структуру</div>
	</form>
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
