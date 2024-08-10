<script setup lang="ts">
import { ElementImage } from '../components/UI/index'

</script>
<template>
    <main class="flex justify-between w-full">
        <div class="p-3 top-players">
            <h2>Топ игроков</h2>
            <ul id="topPlayersList"></ul>
        </div>
        <div class="flex flex-col items-center max-w-[900px] gap-2 p-5 grow-[3] ">
            <!-- There should be a timer -->
            <ElementImage class="h-auto w-full max-w-2xl " elname="H" />
            <button @click="">Поднять руку</button>
        </div>
        <div class="p-3">
            <h2>Поднятые руки</h2>
            <div class="messages">
            </div>
        </div>
    </main>
</template>