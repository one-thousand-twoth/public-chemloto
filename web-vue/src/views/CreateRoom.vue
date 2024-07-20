<script setup lang="ts">
import { RoomInfo, useRoomsStore } from '@/stores/useRoomsStore';
import { ref } from 'vue';

const roomStore = useRoomsStore()

function onSubmit(){
    roomStore.CreateGame(room.value.name)
}

const room = ref<RoomInfo>({
    name: ''
})

</script>

<template>
    <form @submit.prevent="onSubmit()">
        <div class="flex flex-col gap-4">
            <section>
                <label for="roomName">Название комнаты:</label>
                <input v-model="room.name" type="text" name="roomName" required placeholder="Комната 402">
            </section>
            <section>
                <div>
                    <label for="autoPlay">Авто-игра:</label>
                    <input class="m-2" type="checkbox" name="autoPlay">
                </div>
                <!-- Этот блок будет появляться при выборе "Авто-игра" -->
                <section >
                    <label>Время на ход (в секундах):</label>
                    <input type="number" id="timeOptions" name="timeOptions" min="1"
                        placeholder="Введите значение (секунды)">

                </section>
            </section>
            <section>
                <label for="maxPlayers">Макс. количество игроков:</label>
                <input type="number" id="maxPlayers" min="2" placeholder="Например, 24">
            </section>
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