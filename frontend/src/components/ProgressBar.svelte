<script>

  let {
    item = {},
    
    warnL = item.MinY ?? 0,
    warnH = item.MaxY ?? 0,
    alarmL = item.MinR ?? 0,
    alarmH = item.MaxR ?? 0,
    warnLowEnable = item.MinY ? true : false,
    warnHighEnable = item.MaxY ? true : false,
    alarmLowEnable = item.MinR ? true : false,
    alarmHighEnable = item.MaxR ? true : false,

    value = item.value ?? 0
  } = $props();

  let min = $derived(item.min);
  let max = $derived(item.max);
  let percentage = $derived((max - min) === 0 ? 0 : ((value - min) / (max - min)) * 100);

  // реактивный класс для индикатора
  let barColorClass = $derived(
    (alarmLowEnable && value < alarmL) ? 'alarm-low'
    : (alarmHighEnable && value > alarmH) ? 'alarm-high'
    : (warnLowEnable && value < warnL) ? 'warning-low'
    : (warnHighEnable && value > warnH) ? 'warning-high'
    : ''
  );

</script>

<div class="progress-bar-container {barColorClass}" title="{percentage} %">
    <div 
      class="progress-bar-fill {barColorClass}"
      style="width: {percentage}%"
    ></div>
</div>

<style>
  .progress-bar-container {
    --main-color: #23a517;
    --warning-color: #ffd900;
    --alarm-color: #ff0000;

    margin: 6px 0;
    position: absolute;
    inset: 0;
    overflow: hidden;
    /* background-color: #eeeeee; */
    border-radius: 4px;
    border: 1px solid var(--main-color);
  }

  .progress-bar-fill {
    height: 100%;
    background-color: var(--main-color);
    mask-repeat: repeat-x;
  }

  .warning-low, .warning-high {
    border-color: var(--warning-color)
  }
  .alarm-low, .alarm-high {
    border-color: var(--alarm-color)
  }

  .warning-low, .warning-high {
    background-color: var(--warning-color)
  }
  .alarm-low, .alarm-high {
    background-color: var(--alarm-color)
  }

</style>
