<!-- <script setup lang="ts">
import { ref } from 'vue';
import RoomList from '@/components/RoomList.vue';
import ClientList from '@/components/ClientList.vue';
const tabs = {
    "Комнаты": RoomList,
    "Игроки": ClientList,
}

const picked = ref('Комнаты')
</script>
<template>
    <div class="flex flex-wrap shadow-lg">
        <div class="p-6" v-for="(_, tab) in tabs">
            <input class="peer absolute opacity-0" name="tabs" type="radio" :id="tab" :value="tab" v-model="picked" />
            <label
                class="peer-checked:text-black w-auto cursor-pointer bg-slate-50 font-bold text-blue-800 transition-colors hover:bg-white"
                :for="tab">
                {{ tab }}
            </label>
        </div>
    </div>
    <component :is="(<any>tabs)[picked]" class="tab"></component>
</template>

<style>
/* .tabs {
    display: flex;
    flex-wrap: wrap;
    max-width: 700px;
    background: #e5e5e5;
    box-shadow: 0 48px 80px -32px rgba(0, 0, 0, 0.3);
}

.input {
    position: absolute;
    opacity: 0;
} */

/* .label {
    width: 100%;
    padding: 20px 30px;
    background: #e5e5e5;
    cursor: pointer;
    font-weight: bold;
    font-size: 18px;
    color: #7f7f7f;
    transition: background 0.1s, color 0.1s;
} */
/* 
.label:hover {
    background: #d8d8d8;
}

.label:active {
    background: #ccc;
}

.input:focus+.label {
    z-index: 1;
}

.input:checked+.label {
    background: #fff;
    color: #000;
}

@media (min-width: 600px) {
    .label {
        width: auto;
    }
}

.panel {
    display: none;
    padding: 20px 30px 30px;
    background: #fff;
}

@media (min-width: 600px) {
    .panel {
        order: 99;
    }
}

.input:checked+.label+.panel {
    display: block;
} */
</style> -->