<script setup lang="ts">
// import { useGameStore } from '@/stores/useGameStore';

// const GameStore = useGameStore()
// const emit = defineEmits<{
//     selectPlayer: [name: string]
// }>()

</script>
<template>

    

    <!-- <div class=" flex flex-col gap-2
    bg-white bars p-2 h-dvh  md:max-h-[70vh] max-w-[80%]
    " >
        <ul class="list-none flex-initial p-0 font-bold m-0 overflow-y-scroll">
            <li @click="emit('selectPlayer', pl)" class="break-words text-xs flex flex-wrap justify-between items-center p-1 px-2 border-2 border-b-4 border-playing
                    hover:underline 
                    rounded-large my-2"
                v-for=" pl in  Object.keys(Polymers['Альфа'])"
                :key="pl">
                <div class=" inline-flex gap-1 items-center" v-bind="$attrs">
                     <component :is="IconRole(pl.Role)" class="size-4 text-slate-700" /> 
    <span class="text-slate-700"> {{ pl }}</span>
    </div>

     <span>{{ pl.Score }}</span>
    </li>
    </ul>
    <button @click="closeModal()">Закрыть</button>
    </div> -->
</template>