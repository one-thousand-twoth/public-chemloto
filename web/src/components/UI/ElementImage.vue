<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	elname: string
}>()
const imgUrl = computed (() => new URL(`../../assets/items/${props.elname}.svg`, import.meta.url).href)
</script>

<template>
  <div class="relative"><img class="w-full drop-shadow-large h-auto" :src="imgUrl" /></div>
</template>