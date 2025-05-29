<script setup lang="ts">
import { ChemicalElementFormInput } from '@/components/UI/';
import { CreateRoomRequest } from '@/models/RoomModel';
import { useRoomsStore } from '@/stores/useRoomsStore';
import { ref } from 'vue';


const roomStore = useRoomsStore()

function onSubmit() {
    roomStore.CreateGame(room.value)
    emit('exitPanel')
}

const emit = defineEmits(['exitPanel'])


const room = ref<CreateRoomRequest>({
    name: '',
    type: "polymers",
    engineConfig: {
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
        time: 0,
        isAuto: false,
        isAutoCheck: true,
    }

})

function onChangeMaxPlayers(_: Event) {
    if (room.value.type == "polymers") {
        const tradeCount = Math.round(room.value.engineConfig.maxPlayers / 2)
        room.value.engineConfig.elementCounts["TRADE"] = (tradeCount > 10) ? 10 : tradeCount
    }
}

</script>

<template>
    <form class="bg-white" v-if="room.type == 'polymers'" @submit.prevent="onSubmit()">
        <div class="  p-4  relative flex flex-col gap-4">
            <div class="flex">
                <section class="w-3/4">
                    <label for="roomName">Название комнаты:</label>
                    <input class="" v-model="room.name" type="text" name="roomName" required
                        placeholder="Комната 402">
                </section>
            </div>
            <!-- <IconButton class="absolute -right-2 -top-2 p-0 mx-auto text-gray-500" :icon="XMarkIcon" @click="$emit('exit')" /> -->
            <section>
                <div>
                    <label for="autoPlay">Авто-игра:</label>
                    <input v-model="room.engineConfig.isAuto" class="m-2" type="checkbox" name="autoPlay">
                </div>
                <!-- Этот блок будет появляться при выборе "Авто-игра" -->
                <section v-if="room.engineConfig.isAuto">
                    <label>Время на ход (в секундах):</label>
                    <input v-model="room.engineConfig.time" type="number" id="timeOptions" name="timeOptions" min="1"
                        placeholder="Введите значение (секунды)">
                </section>
            </section>
            <section>
                <label for="maxPlayers">Макс. количество игроков:</label>
                <input @change="onChangeMaxPlayers" v-model="room.engineConfig.maxPlayers" type="number" id="maxPlayers"
                    min="2" placeholder="Например, 24">
            </section>
            <details>
                <summary class="mb-4">
                    Количество элементов
                </summary>
                <div class="flex flex-wrap justify-evenly gap-1 ">
                    <ChemicalElementFormInput :max="100" elname="TRADE"
                        v-model="room.engineConfig.elementCounts['TRADE']" />
                    <ChemicalElementFormInput
                        v-for="[name, _] in Object.entries(room.engineConfig.elementCounts).filter(([name, _]) => name != 'TRADE')"
                        :max="100" :elname="name" v-model="room.engineConfig.elementCounts[name]" />
                </div>
            </details>
            <details>
                <summary class="mb-4">
                    Дополнительно
                </summary>
                <label for="isAutoCheck">Проверять игроков:</label>
                <input id='isAutoCheck' v-model="room.engineConfig.isAutoCheck" class="m-2" type="checkbox"
                    name="isAutoCheck">
            </details>

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