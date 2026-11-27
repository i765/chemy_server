<script>
  import SideBar from "./components/SideBar.svelte";
  import AnalogSensors from "./components/AnalogSensors.svelte";
  import DiscreteSensors from "./components/DiscreteSensors.svelte";
  import { config } from "./config.svelte";
  import { sensors } from "./sensors.svelte";
  import { on } from "svelte/events";
  import { onMount } from "svelte";

  let hashParams = $state("");

  function parseHash() {
    // Убираем первый символ '#'
    const rawHash = window.location.hash.slice(1);
    // Превращаем строку в удобный объект
    const params = new URLSearchParams(rawHash);
    hashParams = Object.fromEntries(params.entries());
  }

  onMount(() => {
    parseHash();
  });
</script>

<svelte:window onhashchange={parseHash} />

<div class="spa-layout">
  <SideBar />

  <div class="main-content">
    {#if sensors.isReady}
      {#if hashParams.section === "analog"}
        <AnalogSensors {config} />
      {:else if hashParams.section === "discrete"}
        <DiscreteSensors />
      {/if}
    {:else}
      <div class="loader">
        <div class="informer">
          <div class="warning-icon"></div>
          <div class="text">Нет подключения</div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  @keyframes fadeIn {
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  .loader {
    display: flex;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    justify-content: center;
    align-items: center;
    background-color: #fff;
    z-index: 999;
  }
  .informer {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    opacity: 0;
    transform: translateY(-12px);
    animation: fadeIn 0.5s ease-out 1s forwards;
  }
  .informer .warning-icon {
    width: 100px;
    height: 100px;
    mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 285 300'%3E%3Cpath fill='white' fill-rule='evenodd' d='M117 61 Q130 47 143 47 Q157 47 169 61 L255 222 Q268 244 254 260 Q244 273 229 273 H56 Q41 273 31 260 Q17 244 30 222 Z M133 116 Q128 116 128 122 L133 189 Q133 195 139 195 H148 Q153 195 153 189 L158 122 Q158 116 152 116 Z M143 214 A15 15 0 1 0 143 244 A15 15 0 1 0 143 214 Z'/%3E%3C/svg%3E");
    mask-repeat: no-repeat;
    mask-position: center;
    background-color: #fc0;
  }
  .informer .text {
    font-weight: 500;
    font-size: 16px;
  }

  .spa-layout {
    display: grid;
    grid-template-columns: 280px 1fr;
    height: 100dvh;
    overflow: hidden;
  }

  .main-content {
    padding: 20px 0;
    overflow-y: auto;
  }
</style>
