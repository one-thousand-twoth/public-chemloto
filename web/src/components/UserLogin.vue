<!-- <script setup lang="ts">
import { useUserStore } from "../stores/useUserStore";
import {Modal} from './UI/index'

const userStore = useUserStore();
let usernameInput = ''
</script>

<template>
	<Modal :show="userStore.UserCreds ? false: true">
	
	    <template #header>
	      <h3 class="font-bold">Введите своё имя</h3>
	    </template>
	
	    <template #body>
	      <div class="flex flex-row gap-2">
	        <input type="text" v-model="usernameInput" />
	        <button @click="userStore.Login(usernameInput);">Войти</button>
	      </div>
	    </template>
	
	    <template #footer>
	      <span></span>
	    </template>
	
	  </Modal>
</template> -->
