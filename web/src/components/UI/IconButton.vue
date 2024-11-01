<script setup lang="ts">
import type { FunctionalComponent, HTMLAttributes, VNodeProps } from 'vue';
defineOptions({
  inheritAttrs: false,
});
// defineProps(["icon"]);
defineProps<{
  icon: FunctionalComponent<HTMLAttributes & VNodeProps>
}>()
</script>
<template>
  <div class="flex flex-col" v-bind="$attrs">
    <div class="flex justify-center align-middle transition rounded p-2 m-auto hover:bg-slate-100">
      <component :is="icon" class="size-6 text-slate-500" />
    </div>
  </div>
</template>