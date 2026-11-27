<script>
  import AnalogSensorItem from "./AnalogSensorItem.svelte";
  import { sensors } from "../sensors.svelte";
  import { config } from "../config.svelte";

  
</script>

<div class="div-table">
  {#each config.data.analog as item}
    <AnalogSensorItem {item} value={sensors.analog[item.id]}  status={sensors.warnings[item.id]}/>
  {/each}
</div>

<style>
  .div-table {
    display: grid;
    grid-template-columns: 110px 85px 125px 1fr;
    width: 100%;
  }
</style>
