<script setup lang="ts">
import { IconRole, Role } from '@/models/User';
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
const emit = defineEmits<{
    selectPlayer: [name: string]
}>()

</script>
<template>
    <div>
        <!-- <ul class="list-none p-0 font-bold m-0"> -->
        <TransitionGroup tag="ul" name="fade" class="container list-none p-0 font-bold m-0">
        <li @click="emit('selectPlayer', pl.Name)" class="item cursor-pointer break-words text-xs flex flex-wrap justify-between items-center py-1 px-2 border-2 border-b-4 border-playing
                    hover:underline 
                    rounded-large my-1"
            v-for="pl in GameStore.gameState.Players.filter((pl) => pl.Role === Role.Player).sort((a, b) => b.Score - a.Score)"
            :key="pl.Name">
            <div class=" inline-flex gap-1 items-center" v-bind="$attrs">
                <component :is="IconRole(pl.Role)" class="size-4 text-slate-700" />
                <span class="text-slate-700"> {{ pl.Name }}</span>
            </div>

            <span>{{ pl.Score }}</span>
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

.item {

}

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