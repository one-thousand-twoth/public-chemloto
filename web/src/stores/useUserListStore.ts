import { APISettings } from '@/api/config';
import { Client } from '@/api/core/client';
import { Role, User, UserCreds, UserInfo } from '@/models/User';
import { acceptHMRUpdate, defineStore } from 'pinia';
import { ref } from 'vue';
import { useToasterStore } from "../stores/useToasterStore";
import { useUserStore } from './useUserStore';

export const useUsersListStore = defineStore('usersList', {
  state: () => {
    return {
      UsersList: ref<Array<User>>([]),
      fetching: ref(false)
    }
  },
  actions:
  {
    async fetchUsers() {
      console.log("Fetching users")
      const toasterStore = useToasterStore();
      const userStore = useUserStore();
      this.fetching = true
      if (userStore.UserCreds == null) {
        return
      }
      const client = new Client(APISettings.protocol + APISettings.baseURL, userStore.UserCreds.token);
      try {
        const resp = await client.get("/users");
        if (resp.status == 200) {
          this.UsersList = Object.values(await resp.json())
        }
      } catch (e) {
        toasterStore.error("Не удалось обновить информацию о доступных играх")
      }
      this.fetching = false;
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUsersListStore, import.meta.hot))
}