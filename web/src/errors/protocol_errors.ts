type ErrorResponse = {
  error: string
  form_errors: { [id: string]: string } | null
}

export function assertIsErrorResponse (
  data: unknown
): asserts data is ErrorResponse {
  if (
    typeof data !== 'object' ||
    data === null ||
    !('error' in data) ||
    !('form_errors' in data)
  ) {
    throw new Error("Invalid response: missing 'error' or 'form_error' fields")
  }
}
