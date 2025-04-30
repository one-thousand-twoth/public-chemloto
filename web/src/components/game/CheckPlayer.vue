<script setup lang="ts">
import { WebsocketConnector } from '@/api/websocket/websocket';
import { ChemicalElementFormInput, IconButton } from '@/components/UI';
import { Hand } from '@/models/Game';
import { Field, Polymer, Polymers, StructureNames } from '@/models/Polymers';
import PlayerRoom from '@/views/PlayerRoom.vue';
import {
	XMarkIcon
} from "@heroicons/vue/24/outline";
import { closeModal } from 'jenesius-vue-modal';
import { computed, inject, ref } from 'vue';

const props = defineProps<{
	hand: Hand;
}>()

const ws = inject('connector') as WebsocketConnector


interface CheckStruct<K extends Field> {
	Field: K,
	Name: StructureNames<K> | undefined
}

type AllCheckStructs = {
	[K in Field]: CheckStruct<K>
}[Field];

function Check(Player: string, accept: boolean) {
	ws.Send({
		Type: "ENGINE_ACTION",
		Action: "Check",
		Player: Player,
		Accept: accept,
	})
	closeModal()
}
const check = ref<AllCheckStructs>({
	Field: 'Альфа',
	Name: undefined,
})

const struct = ref<{ [id: string]: number; }>(
	Object.fromEntries(Object.entries(props.hand.Structure).filter(([_, v]) => { return v > 0 }))
)

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
console.log("str", struct.value)
const isHided = ref(true)
const textForBtn = computed(() => { return isHided.value ? 'Показать' : 'Скрыть'})


</script>

<template>
	<form class="z-50 bg-white bars p-4 overflow-y-auto h-[100lvh]  md:max-h-[70vh] max-w-[80%]" @submit.prevent>
		<div class="flex justify-end">

		</div>

		<div class="flex flex-col gap-4 ">
			<div class="flex items-center ">
				<h3 class="text-lg text-center"> Проверка структуры игрока {{ hand.Player.Name }}</h3>
				<IconButton class="text-gray-500" :icon="XMarkIcon" @click="closeModal()" />
			</div>
			<div @click="isHided = !isHided" class="underline cursor-pointer"> {{textForBtn}} </div>
			<div class="text-lg gap-4 flex justify-between">
				<p>Поле: <b>{{ hand.Field }}</b> </p>
				<p>Структура: <b class="inline-block" :class="isHided ? 'bg-main-tint text-main-tint' : '' ">{{ hand.Name }}</b></p>
			</div>

			<div class="flex my-auto flex-wrap justify-between"  :class="isHided ? 'bg-main-tint text-main-tint': '' ">
				<ChemicalElementFormInput :class="isHided ? 'invisible' : ''"  :disabled="true" v-for="elname in Object.keys(struct)" :elname="elname"
					:max="100" v-model.number="struct[elname]" />
			</div>
			<div class="inline-flex gap-1 justify-between">
				<button @click="Check(hand.Player.Name, true)">Принять</button>
				<button class="bg-red-500" @click="Check(hand.Player.Name, false)">Отклонить</button>
			</div>
		</div>
	</form>
</template>
