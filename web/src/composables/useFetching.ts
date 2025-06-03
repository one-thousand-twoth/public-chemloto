import { ref } from 'vue'

export function useFetching<T extends (...args: any[]) => Promise<any>>(fn: T) {
  const isLoading = ref(false)

  const wrapped = async (...args: Parameters<T>): Promise<ReturnType<T>> => {
    isLoading.value = true
    try {
      return await fn(...args)
    } finally {
      isLoading.value = false
    }
  }

  return { isLoading, execute: wrapped }
}