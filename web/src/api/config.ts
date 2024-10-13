const isProd = import.meta.env.MODE === "production"
console.log('prod',isProd)
export const APISettings = {
    token: '',
    headers: new Headers({
        'Accept': 'application/json'
    }),
    baseURL: import.meta.env.PROD ? window.location.hostname+":8090" : window.location.hostname + ':8090',
    protocol: 'http://',
}