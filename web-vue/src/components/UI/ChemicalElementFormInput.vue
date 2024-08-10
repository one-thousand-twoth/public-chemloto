<script setup lang="ts">
const props = defineProps<{
	elname: string;
	value: number
}>()
const emit = defineEmits(['input'])
const url = `../items/${props.elname}.svg`
</script>
<template>
	<div class="flex items-center mb-3">
		<img :src=url :alt=elname class="element-img">

		<input type="number"  min="0" :value="value" @input="$emit('input', value)" />
	</div>
</template>
