<script>
  import SideBar from "./components/SideBar.svelte";
  import AnalogSensors from "./components/AnalogSensors.svelte";
  import { config } from "./config.svelte";
  import { sensors } from "./sensors.svelte";

  $effect(()=>{
    config.fetchNow()

    return sensors.startPolling(1000)
  })
  
</script>

<div class="spa-layout">
  <SideBar />

  <div class="main-content">
    
    <AnalogSensors />

  </div>

</div>
