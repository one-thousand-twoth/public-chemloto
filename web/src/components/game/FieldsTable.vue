<script setup lang="ts">
import { useGameStore } from '@/stores/useGameStore';

const GameStore = useGameStore()
</script>
<template>
    <div class="flex flex-col">
        <div class="-m-1.5 overflow-x-auto">
            <div class="p-1.5 min-w-full inline-block align-middle">
                <div class="overflow-hidden">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead>
                            <tr>
                                <th scope="col" class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
                                    α</th>
                                <th scope="col" class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
                                    β</th>
                                <th scope="col" class="px-6 py-3 text-start text-3xl font-medium text-gray-700 ">
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
</template>
