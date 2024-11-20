<script setup lang="ts">
import type { FunctionalComponent, HTMLAttributes, VNodeProps } from 'vue';
defineOptions({
  inheritAttrs: false,
});
// defineProps(["icon"]);
defineProps<{
  icon: FunctionalComponent<HTMLAttributes & VNodeProps>,
  isChecked: boolean,
}>()
</script>
<template>
  <div class="flex" >
    <div :class="{ 'box': isChecked }" v-bind="$attrs" class="flex border-lg justify-center align-middle transition p-1 m-auto hover:bg-slate-100">
      <component :is="icon" class="size-6 text-slate-500" />
    </div>
  </div>
</template>

<style>

:root {
	--angle: 45deg;
}



.box {
 
/* Size of the border */
--border-size: 2px;
background-origin: border-box;
background-clip:
/* Clip white semi-transparent to the padding-box */
  padding-box, 
  /* Clip colored boxes to the border-box (default) */
  border-box
;

/* Create a dotted border */
border: var(--border-size) dotted transparent;
border-radius: 0.5rem;

/* Create two background layers:
   1. A white semi-transparent
   2. A layer with the colored boxes
 */
background-image:
  linear-gradient(to right, rgb(255 255 255 / 1), rgb(255 255 255 / 1)),

  conic-gradient(
    from var(--angle),
    #ffffff00 0deg 90deg,
    /* #d53e33 0deg 90deg,
    #fbb300 90deg 180deg,
    #377af5 180deg 270deg, */
    #399953 270deg 360deg
  )
;

  animation: rotate 4s linear infinite, opacityChange 3s infinite alternate;

}

@property --angle {
  syntax: "<angle>";
  initial-value: 0deg;
  inherits: false;
}

@keyframes rotate {
  to {
    --angle: 360deg;
  }
}



</style>