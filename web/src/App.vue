<script setup lang="ts">
import { container as ModalContainer } from "jenesius-vue-modal";
import Toaster from './components/Toaster.vue';

defineOptions({
    inheritAttrs: false,
});

</script>
<template>
    <div ref='el'>
        <router-view v-bind="$attrs" />
        <Toaster />
        <modal-container class="z-50"/>
    </div>
</template>