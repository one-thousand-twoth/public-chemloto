import { defineStore } from 'pinia'
import { APISettings } from '@/api/config'
import { ref } from 'vue'
import { Client } from '@/api/core/client'
import { useToasterStore } from "../stores/useToasterStore";

const USER_LOCAL_STORAGE_KEY = 'user'

const getUser = () => {
    const value = localStorage.getItem(USER_LOCAL_STORAGE_KEY)
    return value ? JSON.parse(value) : null;
  }

export const useUserStore = defineStore('users', {
    state: () => {
        return {
            UserCreds:  getUser()  as UserInfo | null
        }
    },
    actions:
     {
        async createUser(input: string) {
            const toasterStore = useToasterStore();
            const client = new Client(APISettings.protocol+APISettings.baseURL, "");
            const resp = await fetch(client.url(`/users`), {
              method: "POST",
              // headers: client.headers(),
              body: new URLSearchParams({
                name: input,
                // IP: "172.16.1.126",
              }),
            });
          
            if (resp.status === 409) {
              toasterStore.error(`Пользователь с именем ${input} уже существует!`);
              return;
            }
          
            if (!resp.ok) {
              console.error("Failed to login with user");
              toasterStore.error(`Не удалось войти под именем ${input}`);
              return;
            }
            
            const token = await resp.text();
            // localStorage.setItem("token", token.value);
            this.UserCreds = {username: input, token: token}
            localStorage.setItem(USER_LOCAL_STORAGE_KEY, JSON.stringify(this.UserCreds));
          
            toasterStore.info(`Вы вошли под пользователем ${this.UserCreds.username}`);
            console.log(`Token for user ${this.UserCreds.username} created: ${this.UserCreds.token}`);
          },
        async check(){
          if (!this.UserCreds){
            return
          }
          const client = new Client(APISettings.protocol+APISettings.baseURL, "");
          try {
            const resp = await client.get("/users/"+this.UserCreds?.token);
            // console.log(resp.status)
            if (resp.status == 200){
            const json = await resp.json();
            // console.log(json)
            this.UserCreds =  json
            } else{
              this.UserCreds = null
            }
          } catch (e) {
            this.UserCreds = null
          }
        },
        $reset(){
          localStorage.removeItem(USER_LOCAL_STORAGE_KEY)
          this.UserCreds = null
        }
    }
})

interface UserInfo {
    username: string
    token: string
    // status: string
}