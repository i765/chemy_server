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

<div class="progress-bar-container {barColorClass}">
    <div 
      class="progress-bar-fill {barColorClass}"
      style="width: {percentage}%"
    ></div>
</div>

<style>
  .progress-bar-container {
    margin: 5px 0;
    position: absolute;
    inset: 0;
    overflow: hidden;
    background-color: #cccccc;
  }

  .progress-bar-fill {
    height: 100%;
    background-color: #31ac0b;
    mask-repeat: repeat-x;
  }

  .progress-bar-container.warning-low, .progress-bar-container.warning-high {
    border-color: #ffd900
  }
  .progress-bar-container.alarm-low, .progress-bar-container.alarm-high {
    border-color: #ff0000
  }

  .progress-bar-fill.warning-low, .progress-bar-fill.warning-high {
    background-color: #ffd900
  }
  .progress-bar-fill.alarm-low, .progress-bar-fill.alarm-high {
    background-color: #ff0000
  }

</style>
