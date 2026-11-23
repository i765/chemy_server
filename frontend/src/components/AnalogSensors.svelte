<script>
  import AnalogSensorItem from "./AnalogSensorItem.svelte";
  import { sensors } from "../sensors.svelte";
  import { config } from "../config.svelte";
  
</script>

<div class="div-table">
  {#each config.data.analog as item}
    <AnalogSensorItem {item} {sensors} />
  {/each}
</div>

<style>
  .div-table {
    display: grid;
    grid-template-columns: 120px 150px 150px 1fr;
    border-top: 1px solid #ccc;
    width: 100%;
  }

</style>
