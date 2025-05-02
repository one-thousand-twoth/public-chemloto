<script setup lang="ts">
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
</script>
<template>
    <div class="flex flex-col">
        <div class="-m-1.5">
            <div class="p-1.5 min-w-full inline-block align-middle">
                <div class="overflow-hidden">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead>
                            <tr>
                                <th scope="col" class="px-3 py-1 text-start  font-medium text-gray-700 ">
                                    α</th>
                                <th scope="col" class="px-3 py-1 text-start  font-medium text-gray-700 ">
                                    β</th>
                                <th scope="col" class="px-3 py-1 text-start  font-medium text-gray-700 ">
                                    γ</th>

                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-200">
                            <tr>
                                <td class="px-3 py-1 whitespace-nowrap   text-gray-800">
                                    {{ GameStore.gameState.Fields["Альфа"]?.Score }}</td>
                                <td class="px-3 py-1 whitespace-nowrap  text-gray-800">{{
                                    GameStore.gameState.Fields["Бета"]?.Score }}</td>
                                <td class="px-3 py-1 whitespace-nowrap  text-gray-800">{{
                                    GameStore.gameState.Fields["Гамма"]?.Score }}</td>

                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>
