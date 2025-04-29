<script setup lang="ts">
import { watch } from 'vue';
import ElementImage from '../UI/ElementImage.vue';
const props = defineProps<{
    // modelValue: number;
    elname: string;
    max: number;
    disabled?: boolean;
}>()

// const emit = defineEmits(['update:modelValue'])
// const url = `../items/${props.elname}.svg`
// const updateValue = (event) => {

// }
const model = defineModel<string>('input_value',{ default: '0' })

function onInput(event: Event) {
  const raw = (event.target as HTMLInputElement).value

  model.value = raw === '' ? '0' : raw
  
}

watch(model, (val) => {
    // @ts-ignore
    if (val === null || val === undefined || val === '') {
        model.value = '0'
    }
    if (Number(val) > 20) {  
        model.value = '20'
    }
})
</script>
<template>
    <div class="relative flex items-center mb-3">
        <ElementImage class="w-10 z-10" :elname="elname" />
        <input  @input="onInput"
            class="relative text-right right-4 w-12  bg-gray-50 border border-gray-500 text-gray-900 font-bold text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block>"
            :disabled="disabled" type="number"  min="0" :value="model" />
    </div>
</template>

<style scoped>
/* Chrome, Safari, Edge, Opera */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* Firefox */
input[type=number] {
  -moz-appearance: textfield;
}
</style>
