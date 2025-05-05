<script setup lang="ts">
import { ElementImage } from '@/components/UI';
import { Player } from '@/models/Game';
import { computed } from 'vue';

const props = defineProps<{
	player: Player
}>()


const bag = computed(()=>{return Object.keys(props.player.Bag).filter((v)=>{return v !== 'TRADE'})})

</script>

<template>
	<div class="" v-if="player">
		<div class="grid grid-cols-2  justify-start my-3 gap-1">
			<div v-for="key in bag" class="flex items-center mb-3 gap-1">
				<ElementImage class="w-7" :elname="key" />
				<div>{{ player.Bag[key] }}</div>
			</div>
		</div>
	</div>
</template>
