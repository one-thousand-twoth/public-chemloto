<script setup lang="ts">
const props = defineProps<{
	elname: string
}>()
const imgUrl = new URL(`../../assets/items/${props.elname}.svg`, import.meta.url).href
console.log(imgUrl)
</script>

<template>
  <img :src="imgUrl" />
</template>