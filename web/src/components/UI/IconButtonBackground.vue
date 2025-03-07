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
  <div class="flex flex-col hover:cursor-pointer">
    <div
      class="flex justify-center align-middle transition rounded p-2 pr-4 m-auto hover:bg-gray-500"
      v-bind="$attrs"
    >
      <component :is="icon" class="size-6 text-white" />
      <slot></slot>
    </div>
  </div>
</template>