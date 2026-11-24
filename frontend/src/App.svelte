<script>
  import SideBar from "./components/SideBar.svelte";
  import AnalogSensors from "./components/AnalogSensors.svelte";
  import { config } from "./config.svelte";
  import { sensors } from "./sensors.svelte";

  let isReady = $state(false);

  $effect(() => {
    if (config.error === null && sensors.error === null) {
      if (config.isLoading && sensors.isLoading) {
        console.log("Загрузка данных...");
      } else {
        console.log("Приложение готово к запуску.");
        isReady = true;
      }
    } else {
      console.error("Ошибка при получении данных!");
      isReady = false;
    }
  });
</script>

<div class="spa-layout">
  
  <SideBar />

  <div class="main-content">
    {#if isReady}
      <AnalogSensors {config} />
    {:else}
      <div class="loader">Подключение к УСО...</div>
    {/if}
  </div>
</div>

<style>
  .loader{
    padding: 30px;
    font-weight: 600;
  }
</style>
