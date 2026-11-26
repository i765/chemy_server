import { mount } from 'svelte'
import './style.css'
import App from './App.svelte'
import { config } from './config.svelte'
import { sensors } from './sensors.svelte'

const app = mount(App, {
  target: document.getElementById('app'),
  sensors: sensors.startPolling() // Запускаем опрос сервера
})

export default app