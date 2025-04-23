<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  modelValue: string | number,
  value: string | number,
  label?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
}>()

const isSelected = computed(() => props.modelValue === props.value)

function select() {
  emit('update:modelValue', props.value)
}
</script>

<template>
  <div
    class="flex items-center translate cursor-pointer px-2 py-2 rounded border"
    :class="isSelected ? 'bg-main-light text-white border-main-dark' : 'text-slate-600 border-slate-300 hover:bg-slate-100'"
    @click="select"
  >
    <!-- <div class="w-4 h-4 rounded-full border-2 flex items-center justify-center"
         :class="isSelected ? 'border-main-dark' : 'border-slate-400'">
      <div v-if="isSelected" class="w-2 h-2 rounded-full bg-main-dark"></div>
    </div> -->
    <slot />
  </div>
</template>
