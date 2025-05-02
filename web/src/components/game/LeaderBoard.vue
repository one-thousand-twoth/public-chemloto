<script setup lang="ts">
import { GameInfo, StateTRADE } from '@/models/Game';
import { IconRole, Role } from '@/models/User';
import { useInterfaceStore } from '@/stores/RoomInterface';
import { useGameStore } from '@/stores/useGameStore';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { ElementImage } from '../UI';

const GameStore = useGameStore()
const { gameState, SelfPlayer } = storeToRefs(GameStore)
const emit = defineEmits<{
  selectPlayer: [name: string]
}>()

const InterfaceStore = useInterfaceStore()
const { currentPlayerSelection } = storeToRefs(InterfaceStore)

const tradeState = computed(() => {
  if (gameState.value.State === "TRADE") {
    return gameState.value as GameInfo & StateTRADE
  }
  return null
})
const players = computed(() => GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score))


const stockList = computed(() => {
  if (!tradeState.value?.StateStruct?.StockExchange.StockList) return [];
  return Object.entries(tradeState.value.StateStruct.StockExchange.StockList);
});

function getStock(username: string) {
  return stockList.value.find((val) => { return val[1].Owner == username })?.[1]
}

function Click(playername: string) {
  currentPlayerSelection.value = playername
  emit('selectPlayer', playername)
}

</script>
<template>
  <div>
    <!-- <ul class="list-none p-0 font-bold m-0"> -->
    <TransitionGroup tag="ul" name="fade" class="container list-none p-0 font-bold m-0">
      <li @click="Click(pl.Name)" class="item cursor-pointer
      break-words text-xs
      items-center py-1 px-2 border-2 border-b-4
      hover:underline 
                    rounded-large my-1" :class="currentPlayerSelection == pl.Name ? 'border-main' : ' border-playing'"
        v-for="pl in players" :key="pl.Name">
        <div class=" flex flex-wrap justify-between">
          <div class=" inline-flex gap-1 items-center" v-bind="$attrs">
            <component :is="IconRole(pl.Role)" class="size-4 text-slate-700" />
            <span class="text-slate-700"> {{ pl.Name }}</span>
          </div>
          <span>{{ pl.Score }}</span>
        </div>
        <div v-if="getStock(pl.Name) !== undefined" class=" text-sm inline-flex gap-1 items-center">
          <span>Отдаёт</span>
          <ElementImage class=" w-8 inline m-1" :elname="getStock(pl.Name)!.Element" />
          <span>за</span>
          <ElementImage class="w-8 inline m-1" :elname="getStock(pl.Name)!.ToElement" />
        </div>
      </li>
    </TransitionGroup>
    <!-- </ul> -->
  </div>
</template>

<style scoped>
.container {
  position: relative;
  padding: 0;
  list-style-type: none;
}

.item {}

/* 1. declare transition */
.fade-move,
.fade-enter-active,
.fade-leave-active {
  transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

/* 2. declare enter from and leave to state */
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: scaleY(0.01) translate(30px, 0);
}

/* 3. ensure leaving items are taken out of layout flow so that moving
      animations can be calculated correctly. */
.fade-leave-active {
  position: absolute;
}
</style>