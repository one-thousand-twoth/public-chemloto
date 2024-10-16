<script setup lang="ts">
const props = defineProps({
  show: Boolean
})
defineEmits(['close'])
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="props.show"
        class="fixed z-50  top-0 left-0 flex justify-center items-center m-auto w-screen h-lvh overflow-y-scroll bg-opacity-20 bg-slate-500">
        <div class="h-fit m-auto bg-white p-4 shadow-lg rounded-sm">
          <div class="">
            <slot name="header"></slot>
          </div>  

          <div class="mt-4">
            <slot name="body">default body</slot>
          </div>

          <div class="mt-4">
            <slot name="footer">
              <button class="modal-default-button" @click="$emit('close')">Закрыть</button>
            </slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
