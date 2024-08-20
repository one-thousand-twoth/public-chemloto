const isProd = import.meta.env.MODE === "production"
console.log('prod',isProd)
export const APISettings = {
    token: '',
    headers: new Headers({
        'Accept': 'application/json'
    }),
    baseURL: import.meta.env.PROD ? window.location.hostname+":1090" : 'localhost:1090',
    protocol: 'http://',
}