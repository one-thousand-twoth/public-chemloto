import { APISettings } from '@/api/config';
import { Client } from '@/api/core/client';
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useToasterStore } from "../stores/useToasterStore";

const USER_LOCAL_STORAGE_KEY = 'user'

const getUser = () => {
  const value = localStorage.getItem(USER_LOCAL_STORAGE_KEY)
  return value ? JSON.parse(value) : null;
}

export const useUserStore = defineStore('users', {
  state: () => {
    return {
      UserCreds: getUser() as UserInfo | null,
      UsersList: ref<Array<UserInfo>>([]),
      fetching: ref(false)
    }
  },
  actions:
  {
    async fetchUsers() {
      console.log("Fetcing users")
      const toasterStore = useToasterStore();
      this.fetching = true
      if (this.UserCreds == null) {
        return
      }
      const client = new Client(APISettings.protocol + APISettings.baseURL, this.UserCreds.token);
      // const token = ref(localStorage.getItem("token") ?? "");
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
    async PatchUser(usr: UserInfo) {
      const toasterStore = useToasterStore();
      const client = new Client(APISettings.protocol + APISettings.baseURL, "");
      let role = ""
      if (usr.role == Role.Player) {
        role = Role.Judge
      } else if (usr.role == Role.Judge) {
        role = Role.Player
      }
      const resp = await fetch(client.url(`/users/${encodeURI(usr.username)}`), {
        method: "POST",
        // headers: client.headers(),
        body: JSON.stringify({
          Role: role
        })
      });
      const json = await resp.json();
      if (!resp.ok) {
        console.error("Failed to login with user");
        toasterStore.error(`Не удалось изменить роль пользователя ${usr.username}`);
        toasterStore.error(json["error"]);
        return;
      }
    },
    async Remove(usr: string) {
      const toasterStore = useToasterStore();
      const client = new Client(APISettings.protocol + APISettings.baseURL, "");

      const resp = await fetch(client.url(`/users/${encodeURI(usr)}`), {
        method: "DELETE",
      });
      const json = await resp.json();
      if (!resp.ok) {
        console.error("Failed to login with user");
        toasterStore.error(`Не удалось удалить пользователя ${usr}`);
        toasterStore.error(json["error"]);
        return;
      }
    },
    async Login(input: string, code: string) {
      const toasterStore = useToasterStore();
      const client = new Client(APISettings.protocol + APISettings.baseURL, "");
      const resp = await fetch(client.url(`/users`), {
        method: "POST",
        // headers: client.headers(),
        body: JSON.stringify({
          name: input,
          code: code,
        })
      });
      if (resp.status === 409) {
        toasterStore.error(`Пользователь с именем ${input} уже существует!`);
        return;
      }

      const json = await resp.json();
      if (!resp.ok) {
        console.error("Failed to login with user");
        toasterStore.error(`Не удалось войти под именем ${input}`);
        json["error"].forEach((element: string) => {
          toasterStore.error(element);
        });
        return;
      }
      // localStorage.setItem("token", token.value);
      this.UserCreds = { username: input, token: json["token"], role: json["role"], room: "" }
      localStorage.setItem(USER_LOCAL_STORAGE_KEY, JSON.stringify(this.UserCreds));

      toasterStore.info(`Вы вошли под именем ${this.UserCreds.username}`);
      console.log(`Token for admin ${this.UserCreds.username} created: ${this.UserCreds.token}`);

    },
    async check() {
      if (!this.UserCreds) {
        return false
      }
      const client = new Client(APISettings.protocol + APISettings.baseURL, "");
      try {
        const token = this.UserCreds.token;
        const resp = await client.get("/users/" + token);
        // console.log(resp.status)
        if (resp.status == 200) {
          const json = await resp.json();
          json["token"] = token
          this.UserCreds = json
          return true
        } else {
          this.UserCreds = null
        }
      } catch (e) {
        this.UserCreds = null
      }
      return false
    },
    $reset() {
      localStorage.removeItem(USER_LOCAL_STORAGE_KEY)
      this.UserCreds = null
    }

  },
  getters:{
    connected: (state) => state.UserCreds?.room ? true : false,
  }
})

export enum Role {
  Admin = "Admin_Role",
  Judge = "Judge_Role",
  Player = "Player_Role",
}
export function i18nRole(role: Role){
  switch (role){
    case Role.Admin:
      return "Админ"
    case Role.Judge:
      return "Судья"
    case Role.Player:
      return "Player"
  }
}
export function emojiRole(role: Role){
  switch (role){
    case Role.Admin:
      return "🅰"
    case Role.Judge:
      return "🛠"
    case Role.Player:
      return "✌"
  }
}
export interface UserInfo {
  username: string
  token: string
  role: Role
  room: string
  // status: string
}