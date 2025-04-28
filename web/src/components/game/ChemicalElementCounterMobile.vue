<script setup lang="ts">
import { useKeyboardStore } from '@/stores/useRaiseHand';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted, useTemplateRef, watch } from 'vue';
import ElementImage from '../UI/ElementImage.vue';

const props = defineProps<{
  elname: string,
  max: number,

  selector: string | number
}>()

const keyboardStore = useKeyboardStore()
const { InputsValues, InputName } = storeToRefs(keyboardStore)

const instance = useTemplateRef<HTMLInputElement>('input_div' + props.elname)


onMounted(() => {
  console.log("mounted " + props.elname)
  if (instance.value === null) {
    console.error("nullish input value")
    return
  }
  // Inputs.value["raise_input_" + props.elname] = instance.value
  InputsValues.value["raise_input_" + props.elname] = instance.value.value // :)
})
onUnmounted(() => {
  // delete Inputs.value["raise_input_" + props.elname]
  delete InputsValues.value["raise_input_" + props.elname]

})

const model = defineModel<number | string>("input_value", { default: 0 })
const selected = defineModel<string | number>("selected")

const isSelected = computed(() => selected.value === props.selector)

function select() {
  console.log(props.selector)
  selected.value = props.selector
  InputName.value = "raise_input_" + props.elname
  instance.value?.focus()
}



watch(InputsValues, (values) => {
  const val = values["raise_input_" + props.elname]
  console.log(val)
  // @ts-ignore
  if (val === null || val === undefined || val === '') {
    model.value = 0
  }
  model.value = val
  if (Number(val) > 20) {
    model.value = 20
    values["raise_input_" + props.elname]= '20'
  }
 
}, {deep: true})
</script>
<template>
  <div :class="isSelected ? 'drop-shadow-large' : ''" @click.stop="select()" class="relative flex items-center mb-3">
    <ElementImage class="w-10 z-10" :elname="elname" />

    <!-- <input @input="onInput" :class="isSelected ? 'bars' : ''"
      class="relative text-right right-4 w-12  bg-gray-50 border border-gray-500 text-gray-900 font-bold text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block>"
      type="number" disabled :max="max" min="0" :value="model" /> -->
    <input readonly @focus="select()" :ref="'input_div' + elname" class="relative text-right right-4 w-12  bg-gray-50 border border-gray-500 text-gray-900 font-bold text-sm rounded-lg
      focus:border-main block  px-2 py-1 " :class="isSelected ? 'border-2 border-main' : ''" type="number" :max="max" min="0" v-model="model">
    </input>

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
