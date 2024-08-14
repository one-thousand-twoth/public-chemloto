<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	elname: string
}>()
const imgUrl = computed (() => new URL(`../../assets/items/${props.elname}.svg`, import.meta.url).href)
console.log(imgUrl)
</script>

<template>
  <img :src="imgUrl" />
</template>