<template>
    <div @click.stop :class="keyboardClass"></div>
</template>

<script>
import { useKeyboardStore } from "@/stores/useRaiseHand";
import { mapStores, mapWritableState } from "pinia";
import Keyboard from "simple-keyboard";
import "simple-keyboard/build/css/index.css";


export default {
    name: "SimpleKeyboard",
    computed: {
        ...mapWritableState(useKeyboardStore, ['Value', 'InputsValues', 'InputName', 'Keys'])
    },
    props: {
        keyboardClass: {
            default: "simple-keyboard",
            type: String
        },
        input: {
            type: String
        }
    },
    data: () => ({
        keyboard: null
    }),
    mounted() {
        this.keyboard = new Keyboard(this.keyboardClass, {
            onChange: this.onChange,
            onKeyPress: this.onKeyPress,
            inputName: this.InputName,
            layout: {
                default: [
                    "1 2 3",
                    "4 5 6",
                    "7 8 9",
                    "0 {bksp}",
                ]
            },
        });
    },
    methods: {
        onChange(input) {
            this.$emit("onChange", input);
            console.log("onChange", input);
            // this.Value = input

            this.InputsValues[this.InputName] = input.replace(/^0+/, '');
        },
        onKeyPress(button) {
            this.$emit("onKeyPress", button);
            console.log("onKeyPress", button);
            /**
             * If you want to handle the shift and caps lock buttons
             */
            if (button === "{shift}" || button === "{lock}") this.handleShift();
        },
    },
    watch: {
        InputName(InputName) {
            console.log(
                'SimpleKeyboard: inputName updated',
                InputName,
            );
            this.keyboard.setOptions({ inputName: InputName });
        },
        InputsValues: {
            handler(InputsValues) {
                console.log(
                    'SimpleKeyboard: inputs Updated',
                    this.keyboard.options.inputName,
                    InputsValues
                );
                this.keyboard.replaceInput(InputsValues);
            },
            deep: 1,
        },
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>