// Types for the result object with discriminated union
type Success<T> = {
  data: T
  error: null
}

type Failure<E> = {
  data: null
  error: E
}

export type Result<T, E = Error> = Success<T> | Failure<E>

// Main wrapper function
export async function tryCatch<T, E = Error> (
  promise: Promise<T>
): Promise<Result<T, E>> {
  try {
    const data = await promise
    return { data, error: null }
  } catch (error) {
    return { data: null, error: error as E }
  }
}

export class AppError extends Error{
  
}

export class FormValidationError extends Error {
  public readonly fields: Record<string, string>

  constructor (fields: Record<string, string>, message?: string) {
    super(message ?? 'Ошибка валидации формы')
    this.name = 'FormValidationError'
    this.fields = fields

    Object.setPrototypeOf(this, FormValidationError.prototype)
  }
}
