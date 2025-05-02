import { acceptHMRUpdate, defineStore } from 'pinia';



export const useInterfaceStore = defineStore('interface', {
  state: (): {
    currentPlayerSelection: string | undefined
  } => ({
      currentPlayerSelection: undefined
  }),
  actions: {


  },
  getters:{
    // Keys:  (state) => Object.keys(state.InputsValues)
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useInterfaceStore, import.meta.hot))
}
