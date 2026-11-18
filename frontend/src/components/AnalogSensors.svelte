<script>
  import AnalogSensorItem from "./AnalogSensorItem.svelte"
  import { sensors } from "../sensors.svelte"
  import { config } from "../config.svelte";

  $effect(() => {

    console.log(config.items.analog_channels)
    
  })

</script>

<div class="div-table">
  <div class="div-cell head">Название</div>
  <div class="div-cell head">Шкала</div>
  <div class="div-cell head">Значение</div>
  <div class="div-cell head">Описание</div>

  <AnalogSensorItem itemName="asd" />


  {#if sensors.isLoading && sensors.length === 0}
    <p>Первичная загрузка...</p>
  {:else if sensors.error}
    <p class="error">Ошибка связи: {sensors.error}</p>
  {:else}

  <p>{sensors.items.analogValues}</p>



  {/if}
</div>

<style>
  .div-table {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    border-top: 1px solid #333;
    border-left: 1px solid #333;
    width: 100%;
    font-family: sans-serif;
  }

  .div-cell {
    padding: 2px 10px;
    border-right: 1px solid #333;
    border-bottom: 1px solid #333;
  }
</style>
