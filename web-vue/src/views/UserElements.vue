<script setup lang="ts">
import { Player } from '@/stores/useGameStore';
import { ElementImage } from '@/components/UI';

defineProps<{
	player: Player 
}>()

</script>

<template>
	<div class="" v-if="player">
		<span>Роль: {{ player.Role }}</span>
		<div class="flex flex-wrap  justify-start my-3 gap-3">
			<div v-for="key in Object.keys(player.Bag)" class="flex items-center mb-3 gap-1">
				<ElementImage class="w-8" :elname="key" />
				<div>{{ player.Bag[key] }}</div>
			</div>
		</div>
	</div>
</template>
