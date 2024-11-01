<script setup lang="ts">
import { IconRole, UserInfo, Role } from '@/models/User';
// defineProps(["icon"]);
defineProps<{
    role: Role,
    name: string,
}>()
</script>
<template>
    <div class=" inline-flex gap-1" v-bind="$attrs">
        <component :is="IconRole(role)" class="size-6 text-slate-700" />
        <span class="text-slate-700"> {{ name }}</span>
    </div>
</template>