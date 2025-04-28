import { acceptHMRUpdate, defineStore } from 'pinia';



// const struct = ref<{ [id: string]: number; }>(
// 	Object.fromEntries(Object.entries(props.player.Bag).map(([name]) => { return [name, 0] }))
// )

export const useKeyboardStore = defineStore('value', {
  state: (): {
    Value: string
    // Inputs: {[id: string]: HTMLInputElement}
    InputsValues: {[id: string]: string}
    InputName: string
  } => ({
      Value: '',
      // Inputs: {},
      InputName: '',
      InputsValues: {},
  }),
  actions: {


  },
  getters:{
    Keys:  (state) => Object.keys(state.InputsValues)
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKeyboardStore, import.meta.hot))
}
