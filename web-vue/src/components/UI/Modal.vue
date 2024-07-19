<!-- <template>
    <div class="absolute flex justify-center items-center h-screen w-screen bg-opacity-20 bg-slate-500">
      <div class="h-fit m-auto bg-white p-4 shadow-lg rounded-lg">
        <slot />
      </div>
    </div>
</template> -->

<script setup lang="ts">
const props = defineProps({
  show: Boolean
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
    <div v-if="props.show"
      class="fixed -z-50 top-0 left-0 flex justify-center items-center m-auto w-screen h-screen bg-opacity-20 bg-slate-500">
      <div class="h-fit m-auto bg-white p-4 shadow-lg rounded-sm">
        <div class="modal-header">
          <slot name="header">default header</slot>
        </div>

        <div class="modal-body mt-4">
          <slot name="body">default body</slot>
        </div>

        <div class="modal-footer">
          <slot name="footer">
            <button class="modal-default-button" @click="$emit('close')">Закрыть</button>
          </slot>
        </div>
      </div>
    </div>
    </Transition>
  </Teleport>
</template>

<!-- <style>
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  transition: opacity 0.3s ease;
}

.modal-container {
  width: 300px;
  margin: auto;
  padding: 20px 30px;
  background-color: #fff;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-header h3 {
  margin-top: 0;
  color: #42b983;
}

.modal-body {
  margin: 20px 0;
}

.modal-default-button {
  float: right;
}
-->
<style>
/*
 * The following styles are auto-applied to elements with
 * transition="modal" when their visibility is toggled
 * by Vue.js.
 *
 * You can easily play with the modal transition by editing
 * these styles.
 */

.modal-enter-from {
  opacity: 0;
}

.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  -webkit-transform: scale(1.1);
  transform: scale(1.1);
}
</style>