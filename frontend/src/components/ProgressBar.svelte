<script>

  let {
    item = {},
    warning = false,
    alarm = false,
    value = 0,
    status = 0
  } = $props();

  let min = $derived(item.min);
  let max = $derived(item.max);

  let percentage = $derived((max - min) === 0 ? 0 : ((value - min) / (max - min)) * 100);

  

</script>

<div class="progress-bar-container status-{status}" title="{percentage} %">
    <div 
      class="progress-bar-fill status-{status}"
      style="width: {percentage}%"
    ></div>
</div>

<style>
  .progress-bar-container {
    --main-color: #45e252;
    --warning-color: #ffc013;
    --alarm-color: #ff5050;
    margin: 6px 0;
    position: absolute;
    inset: 0;
    overflow: hidden;
    border-radius: 3px;
    border: 1px solid var(--main-color);
  }

  .progress-bar-fill {
    height: 100%;
    background-color: var(--main-color);
  }

  .progress-bar-container.status-1{
    border-color: var(--warning-color);
  }
  .progress-bar-fill.status-1 {
    background-color: var(--warning-color)
  }
  .progress-bar-container.status-2{
    border-color: var(--alarm-color);
  }
  .progress-bar-fill.status-2{
    background-color: var(--alarm-color)
  }

</style>
