<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	elname: string
}>()
const imgUrl = computed (() => new URL(`../../assets/items/${props.elname}.svg`, import.meta.url).href)
</script>

<template>
  <div class="relative"><img class=" drop-shadow-large h-full w-auto aspect-square" :src="imgUrl" /></div>
</template>