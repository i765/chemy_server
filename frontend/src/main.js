import { mount } from 'svelte'
import './style.css'
import App from './App.svelte'
import { config } from './config.svelte'
import { sensors } from './sensors.svelte'

const app = mount(App, {
  target: document.getElementById('app'),
  config: config.fetchNow(),
  sensors: sensors.startPolling(1000)
})

export default app