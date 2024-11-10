<script setup lang="ts">
import Toaster from './components/Toaster.vue';

defineOptions({
    inheritAttrs: false,
});



</script>
<template>
    <router-view v-bind="$attrs" />
    <Toaster />
</template>