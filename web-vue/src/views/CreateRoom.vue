<script lang="ts">
import { defineComponent, computed, ref } from "vue"

export default defineComponent({})


</script>

<template>
    <label for="roomName">Название комнаты:</label>
    <input type="text" name="roomName" id="roomName" required placeholder="Комната 402">
    <label for="autoPlay">Авто-игра:</label>
    <input type="checkbox" name="autoPlay" id="autoPlay">

    <!-- Этот блок будет появляться при выборе "Авто-игра" -->
    <div id="timeControl" style="display: none;">
        <label>Время на ход (в секундах):</label>
        <input type="number" id="timeOptions" name="timeOptions" min="1" placeholder="Введите значение (секунды)">
    </div>
    <label for="maxPlayers">Макс. количество игроков:</label>
    <input type="number" id="maxPlayers" min="2" placeholder="Например, 24">
</template>
