<script setup lang="ts">
import ElementImage from './ElementImage.vue';
defineProps<{
	// modelValue: number;
	elname: string;
	max: number;
}>()
// const emit = defineEmits(['update:modelValue'])
// const url = `../items/${props.elname}.svg`
// const updateValue = (event) => {
    
// }
const model = defineModel()
</script>
<template>
	<div class="flex items-center mb-3 gap-3">
		<ElementImage class="w-8" :elname="elname" />

		<input type="number" :max="max" min="0" v-model="model" />
	</div>
</template>
