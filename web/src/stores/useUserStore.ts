import { APISettings } from '@/api/config';
import { Client } from '@/api/core/client';
import { Role, User, UserCreds, UserInfo } from '@/models/User';
import { acceptHMRUpdate, defineStore } from 'pinia';
import { useToasterStore } from "../stores/useToasterStore";

const USER_LOCAL_STORAGE_KEY = 'user'

const getUser = () => {
  const value = localStorage.getItem(USER_LOCAL_STORAGE_KEY)
  return value ? JSON.parse(value) : null;
}

export const useUserStore = defineStore('users', {
  state: () => {
    return {
      // изменение UserCreds вызывает обновление вебсокета
      UserCreds: getUser() as UserCreds | null,
      UserInfo: { room: "", role: Role.Player } as UserInfo,
      fetching: false
    }
  },
  actions:
  {
    getUser() {
      return { ...this.UserCreds, ...this.UserInfo }
    },

    async PatchUser(usr: User) {
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
    async Login(input: string, code: string): Promise<{
      success: boolean;
      formErrors?: Record<string, string>;
    }> {
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
        return { success: false };
      }

      const json = await resp.json();
      if (!resp.ok) {
        console.error("Failed to login with user");

         // Проверяем есть ли ошибки формы
        if (json["form_errors"]) {
          // Возвращаем ошибки формы в компонент
          return { 
            success: false, 
            formErrors: json["form_errors"] 
          };
        }

        toasterStore.error(`Не удалось войти под именем ${input}`);
        json["error"].forEach((element: string) => {
          toasterStore.error(element);
        });
        return { success: false };
      }
      this.UserCreds = { username: input, token: json["token"], }
      this.UserInfo = { room: json["room"] ?? "", role: json["role"] }
      localStorage.setItem(USER_LOCAL_STORAGE_KEY, JSON.stringify(this.UserCreds));

      toasterStore.info(`Вы вошли под именем ${this.UserCreds.username}`);
      console.log(`Token for admin ${this.UserCreds.username} created: ${this.UserCreds.token}`);
      return { success: true };
    },
    async check() {
      let self = this
      this.fetching = true
      const ok = await (async function () {
        // do something

        if (!self.UserCreds) {
          return false
        }

        const client = new Client(APISettings.protocol + APISettings.baseURL, "");
        try {
          const token = self.UserCreds.token;
          const resp = await client.get("/users/" + token);
          // console.log(resp.status)
          if (resp.status == 200) {
            const json = await resp.json();
            // json["token"] = token
            self.UserCreds = { username: json["username"], token: token }
            self.UserInfo = { room: json["room"], role: json["role"] }
            return true
          } else {
            self.UserCreds = null
          }
        } catch (e) {
          self.UserCreds = null
        }
        return false
      })()
      this.fetching = false
      return ok
    },
    $reset() {
      localStorage.removeItem(USER_LOCAL_STORAGE_KEY)
      this.UserCreds = null
    }

  },
  getters: {
    connected: (state) => state.UserInfo?.room != "",
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
}