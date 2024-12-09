import { acceptHMRUpdate, defineStore } from "pinia";

// Status will define toast color and icon
export type TToastStatus = "success" | "info" | "error";

const defaultTimeout: number = 4000;

// timeout is conditional because we will define default one
type ToastPayload = { timeout?: number; text: string };

const createToast = (text: string, status: TToastStatus): IToast => ({
    text,
    status,
    id: Math.random() * 1000,
});

interface IToast {
    // Text of toast
    text: string;
    status: TToastStatus;
    // Id to differentiate toasts
    id: number;
}

 

export const useToasterStore = defineStore("toaster-store", {
    state: (): { toasts: IToast[], callback: (() => void) | null } => ({
        toasts: [],
        callback: null
    }),
    actions: {
        updateState(payload: ToastPayload, status: TToastStatus) {
            if (this.callback != null){
                this.callback()
            }
            // Get text and timeout from payload
            const { text, timeout } = payload;
            // We create the toast with function above
            const toast = createToast(text, status);

            // We push toasts to the state
            this.toasts.push(toast);

            // We create a delay to delete toast after its provided timeout is over
            setTimeout(() => {
                this.toasts = this.toasts.filter((t) => t.id !== toast.id);
            }, timeout ?? defaultTimeout);
        },
        success(payload: string) {
            this.updateState({ text: payload }, "success");
        },

        info(payload: string) {
            this.updateState({ text: payload }, "info");
        },

        error(payload: string) {
            this.updateState({ text: payload }, "error");
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useToasterStore, import.meta.hot))
}
