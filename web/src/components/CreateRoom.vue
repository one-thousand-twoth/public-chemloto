<script setup lang="ts">
import ChemicalElementFormInput from '@/components/UI/ChemicalElementFormInput.vue';
import { RoomInfo, useRoomsStore } from '@/stores/useRoomsStore';
import { ref } from 'vue';

const roomStore = useRoomsStore()

function onSubmit() {
    roomStore.CreateGame(room.value)
}

const room = ref<RoomInfo>({
    name: '',
    maxPlayers: 0,
    elementCounts: {
        "TRADE": 4,
        "O": 28,
        "N": 16,
        "H": 52,
        "Cl": 16,
        "CH3": 28,
        "CH2": 24,
        "CH": 24,
        "C6H5": 16,
        "C6H4": 16,
        "C": 40
    },
    engine: {Status: ""},
    time: 0,
    isAuto: false
})

</script>

<template>
    <form @submit.prevent="onSubmit()">
        <div class="flex flex-col gap-4 ">
            <section>
                <label for="roomName">Название комнаты:</label>
                <input v-model="room.name" type="text" name="roomName" required placeholder="Комната 402">
            </section>
            <section>
                <div>
                    <label for="autoPlay">Авто-игра:</label>
                    <input v-model="room.isAuto" class="m-2" type="checkbox" name="autoPlay">
                </div>
                <!-- Этот блок будет появляться при выборе "Авто-игра" -->
                <section v-if="room.isAuto">
                    <label>Время на ход (в секундах):</label>
                    <input v-model="room.time" type="number" id="timeOptions" name="timeOptions" min="1"
                        placeholder="Введите значение (секунды)">
                </section>
            </section>
            <section>
                <label for="maxPlayers">Макс. количество игроков:</label>
                <input v-model="room.maxPlayers" type="number" id="maxPlayers" min="2" placeholder="Например, 24">
            </section>
            <div class="flex flex-wrap justify-evenly gap-1 " >
                <ChemicalElementFormInput v-for="[name, _] in Object.entries(room.elementCounts)" :max="100" :elname="name"  v-model="room.elementCounts[name]" />
            </div>
            <button type="submit">Создать</button>
        </div>
    </form>
</template>

<style scoped>
section {
    display: flex;
    flex-direction: column;
    gap: 0.25rem
}
</style>