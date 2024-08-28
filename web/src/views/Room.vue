<script setup lang="ts">
import { computed, inject, ref, } from 'vue'
import { ElementImage } from '../components/UI/index'
import { Hand, useGameStore } from '@/stores/useGameStore'
import { WebsocketConnector } from '@/api/websocket/websocket'
import { Role, useUserStore } from '@/stores/useUserStore';
import { Modal } from '../components/UI/index';
import UserElements from './UserElements.vue'
import Timer from '@/components/UI/Timer.vue';
import CheckPlayer from './CheckPlayer.vue'
import RaiseHandComp from './RaiseHandComp.vue';
import Trade from './Trade.vue'
import {
    ArrowLeftStartOnRectangleIcon
} from "@heroicons/vue/24/outline";
import IconButtonBackground from '@/components/UI/IconButtonBackground.vue';
const GameStore = useGameStore()
const userStore = useUserStore()
// const player = GameStore.SelfPlayer
console.log(GameStore.gameState.Players)
console.log(userStore.UserCreds?.username)
const ws = inject('connector') as WebsocketConnector
function StartGame() {
    ws.Send({
        Type: 'HUB_STARTGAME',
        Name: GameStore.name.toString()
    })
}

function GetElement() {
    console.log('Get element!')
    ws.Send({
        Type: 'ENGINE_ACTION',
        Action: 'GetElement'
    })
}
function SendContinue() {
    ws.Send({
        Type: 'ENGINE_ACTION',
        Action: 'Continue'
    })
}
function DisconnectGame() {
      ws.Send(
        {
          "Type": "HUB_UNSUBSCRIBE",
          "Target": "room",
          "Name": GameStore.name
        }
      )
    }
// function RaiseHand() {
//     console.log('Get element!')
//     ws.Send({
//         Type: 'ENGINE_ACTION',
//         Action: 'RaiseHand'
//     })
// }

const currPlayer = computed(() => {
    return GameStore.gameState.Players.find(player => player.Name === curInfoPlayer.value)
})

// const currCheckPlayer = computed(() => {

//     // return GameStore.gameState.Players.find(player => player.Name === curCheckPlayer.value)
// })

const curInfoPlayer = ref('')
const curCheckPlayer = ref<Hand>()
const RaiseHandButton = ref(false)
const TradeButton = ref(false)


</script>
<template>
    <div class="relative flex max-h-lvh flex-col items-center overflow-x-hidden">
        <main class="flex justify-between w-lvw grow gap-10 bg-gray-100">
            <div class="flex flex-col m-3 w-[20%] gap-2">
                <div class="bars p-3 min-w-[8.5rem]  grow-[1] bg-gray-50">
                    <h2>Топ игроков</h2>
                    <ul class="list-none p-0 font-bold m-0">
                        <li @click="curInfoPlayer = pl.Name" class="break-words flex justify-between items-center p-2
                    hover:underline 
                    [&:nth-child(1)]:bg-amber-300
                    [&:nth-child(2)]:bg-stone-300
                    [&:nth-child(3)]:bg-yellow-500
                    
                    rounded-md my-2"
                            v-for="pl in GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score) "
                            :key="pl.Name">
                            {{ pl.Name }} - {{ pl.Score }}
                        </li>
                    </ul>
                </div>
                <IconButtonBackground v-if="!GameStore.gameState.Started || userStore.UserCreds?.role != Role.Player "class="w-full bg-red-700 text-white  rounded-lg"
                    :icon="ArrowLeftStartOnRectangleIcon" @click="DisconnectGame()">Выйти</IconButtonBackground>
            </div>
            <div class="flex flex-col items-center  h-lvh max-w-[900px] gap-2 p-5 pb-60 grow-[3]">
                <Timer />

                <ElementImage class="h-auto w-full max-w-[50lvh]" :elname="GameStore.currElement" />
                <div class="flex flex-col">
                    <div class="-m-1.5 overflow-x-auto">
                        <div class="p-1.5 min-w-full inline-block align-middle">
                            <div class="overflow-hidden">
                                <table class="min-w-full divide-y divide-gray-200">
                                    <thead>
                                        <tr>
                                            <th scope="col"
                                                class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
                                                α</th>
                                            <th scope="col"
                                                class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
                                                β</th>
                                            <th scope="col"
                                                class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
                                                γ</th>

                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-gray-200">
                                        <tr>
                                            <td class="px-6 py-4 whitespace-nowrap text-3xl font-3xl text-gray-800">
                                                {{ GameStore.gameState.Fields["Альфа"]?.Score }}</td>
                                            <td class="px-6 py-4 whitespace-nowrap text-3xl text-gray-800">{{
                                                GameStore.gameState.Fields["Бета"]?.Score }}</td>
                                            <td class="px-6 py-4 whitespace-nowrap text-3xl text-gray-800">{{
                                                GameStore.gameState.Fields["Гамма"]?.Score }}</td>

                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
                <template v-if="!GameStore.gameState.Started">
                    <button v-if="userStore.UserCreds?.role != Role.Player" @click="StartGame()">
                        Начать игру
                    </button>
                    <button disabled v-else>Ждем
                        начала</button>
                </template>
                <template v-else>
                    <template v-if="GameStore.gameState.State == 'OBTAIN'">
                        <button v-if="userStore.UserCreds?.role == Role.Player"
                            @click="RaiseHandButton = !RaiseHandButton">Поднять
                            руку</button>
                        <button v-if="userStore.UserCreds?.role != Role.Player" @click="GetElement()">Достать
                            элемент</button>
                    </template>
                    <template v-if="GameStore.gameState.State == 'HAND'">
                        <button v-if="userStore.UserCreds?.role == Role.Player"
                            @click="RaiseHandButton = !RaiseHandButton">Поднять
                            руку</button>
                        <button disabled v-if="userStore.UserCreds?.role != Role.Player" @click="GetElement()">Ждем
                            проверки</button>
                    </template>
                    <template v-if="GameStore.gameState.State == 'TRADE' && userStore.UserCreds?.role != Role.Player">
                        <button @click="TradeButton = !TradeButton">Обменять</button>
                        <button @click="SendContinue()">Продолжить</button>
                    </template>
                </template>
                <!-- <div class="mb-12" id="empty"></div> -->
            </div>
            <div class="bars min-w-[135px] p-3 bg-gray-50 w-[20%] m-3">
                <h2>Поднятые руки</h2>
                <ul class="list-none p-0 font-bold m-0">
                    <li @click="curCheckPlayer = pl" class="break-words flex justify-between items-center p-2 hover:underline rounded-md my-2 mx-0
                    border-solid border-2 border-gray-600 m-3" v-for="pl in GameStore.gameState.RaisedHands">
                        {{ pl.Player.Name }} - {{ pl.Field }}
                    </li>
                </ul>
            </div>
        </main>
        <div class="fixed bottom-3 right-3 flex flex-wrap items-center" id="lastElementsContainer">
            <h3 class="mr-2">Последние элементы</h3>
            <ElementImage v-for="el in GameStore.LastElements" :elname="el" />
        </div>
        <Modal :show="currPlayer !== undefined" @close="curInfoPlayer = ''">
            <template #header>
                <h3 class="font-bold text-center">Информация о игроке {{ curInfoPlayer }}</h3>
            </template>
            <template #body>
                <UserElements v-if="currPlayer" :player="currPlayer" />
            </template>
        </Modal>
        <Modal v-if="userStore.UserCreds?.role != Role.Player" :show="curCheckPlayer !== undefined"
            @close="curCheckPlayer = undefined">
            <template #header>
                <h3 class="font-bold text-center">Проверка структуры {{ curCheckPlayer?.Player.Name }}</h3>
            </template>
            <template v-if="curCheckPlayer !== undefined" #body>
                <CheckPlayer :player="curCheckPlayer" />
            </template>
        </Modal>
        <Modal v-if="userStore.UserCreds?.role === Role.Player" :show="RaiseHandButton"
            @close="RaiseHandButton = false">
            <template #header>
                <h3 class="font-bold text-center">Поднять руку</h3>
            </template>
            <template v-if="GameStore.SelfPlayer" #body>
                <RaiseHandComp :player="GameStore.SelfPlayer" />
            </template>
        </Modal>
        <Modal v-if="userStore.UserCreds?.role !== Role.Player" :show="TradeButton" @close="TradeButton = false">
            <template #header>
                <h3 class="font-bold text-center">Обменять</h3>
            </template>
            <template v-if="GameStore.SelfPlayer" #body>
                <Trade :players="GameStore.gameState.Players" />
            </template>
        </Modal>
    </div>
</template>

<style scoped>
body {
    background-color: #3b4e5e;
}
</style>
