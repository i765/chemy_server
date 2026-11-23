<script>
  import ProgressBar from "./ProgressBar.svelte";
  let { item, sensors } = $props();

  let value =  $derived(sensors.analog[item.id])
  
</script>

<div class="div-cell id">{item.name}</div>
<div class="div-cell bar"><ProgressBar {item}{value}/></div>
<div class="div-cell">{value}</div>
<div class="div-cell">{item.desc}</div>

<style>
  .div-cell {
    border-bottom: 1px solid #aaaaaa;
  }
  .bar{
    position: relative;
  }
</style>
