<script>

  let {
    item = {},
    warning = false,
    alarm = false,
    value = 0
  } = $props();

  let min = $derived(item.min);
  let max = $derived(item.max);

  let percentage = $derived((max - min) === 0 ? 0 : ((value - min) / (max - min)) * 100);

  let status = $derived(warning ? (alarm ? 'alarm' : 'warning') : '')

</script>

<div class="progress-bar-container {status}" title="{percentage} %">
    <div 
      class="progress-bar-fill {status}"
      style="width: {percentage}%"
    ></div>
</div>

<style>
  .progress-bar-container {
    --main-color: #999999;
    --warning-color: #ffd900;
    --alarm-color: #ff0000;
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

  .progress-bar-container.warning{
    border-color: var(--warning-color);
  }
  .progress-bar-fill.warning {
    background-color: var(--warning-color)
  }
  .progress-bar-container.alarm{
    border-color: var(--alarm-color);
  }
  .progress-bar-fill.alarm {
    background-color: var(--alarm-color)
  }

</style>
