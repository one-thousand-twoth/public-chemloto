<script setup lang="ts">
import { inject } from 'vue';
import { ElementImage } from '../components/UI/index'
import { useGameStore } from '@/stores/useGameStore'
import { WebsocketConnector } from '@/api/websocket/websocket';

const GameStore = useGameStore()

const ws = inject('connector') as WebsocketConnector
function StartGame() {
    ws.Send(
        {
            "Type": "HUB_STARTGAME",
            "Name": GameStore.name
        }
    )
}

function GetElement() {
    console.log("Get element!")
    ws.Send(
        {
            "Type": "ENGINE_ACTION",
            "Action": "GetElement"
        }
    )
}

</script>
<template>
    <div class="relative flex min-h-lvh 
        flex-col items-center overflow-x-hidden 
    ">
        <main class="flex justify-between w-lvw grow gap-20 bg-gray-100">
            <div class="bars p-3 w-[20%] bg-gray-50">
                <h2>Топ игроков</h2>
                <ul class="list-none p-0 font-bold m-0">
                    <li class="break-words flex justify-between items-center p-2
                    hover:underline
                [&:nth-child(1)]:bg-amber-300
                [&:nth-child(2)]:bg-stone-300
                [&:nth-child(3)]:bg-yellow-500
                rounded-md
                my-2
                " v-for="pl in GameStore.gameState.Players"> {{ pl }}</li>
                </ul>
            </div>
            <div class="flex flex-col items-center max-w-[900px] gap-2 p-5 grow-[3]">
                <!-- There should be a timer -->
                <ElementImage class="h-auto w-full max-w-2xl " :elname="GameStore.currElement" />
                <button @click="">Поднять руку</button>
                <button @click="GetElement()">Достать элемент</button>
                <button @click="StartGame()">Начать игру</button>
            </div>
            <div class="bars p-3 bg-gray-50 w-[20%]">
                <h2>Поднятые руки</h2>
                <div class="messages">
                </div>
            </div>
        </main>
        <div class="fixed bottom-3 right-3 flex flex-wrap items-center" id="lastElementsContainer">
            <h3 class="mr-2">Последние элементы</h3>
            <ElementImage v-for="el in GameStore.LastElements" :elname="el" />
        </div>
    </div>
</template>

<style scoped>
body {
    background-color: #3b4e5e;
}
</style>