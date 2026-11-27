<script>
  import ProgressBar from "./ProgressBar.svelte";

  let { item = {}, value = 0, status = 0 } = $props();

  let selected = $state(false);
</script>

<div class="grid-row" class:selected onclick={() => (selected = !selected)}>
  <div class="div-cell name">{item.name}</div>
  <div class="div-cell bar"><ProgressBar {item} {value} {status} /></div>
  <div class="div-cell value" title="{item.min} - {item.max} {item.unit}">
    {value}
    {item.unit}
  </div>
  <div class="div-cell desc">
    {item.desc} ({item.min} - {item.max}
    {item.unit})
  </div>
</div>

<style>
  .grid-row {
    display: grid;
    grid-column: 1 / -1; /* Растягиваем строку на все 4 колонки */
    grid-template-columns: subgrid; /* Наследуем размеры колонок от .grid-table */
    border-radius: 6px;
    cursor: pointer;
    margin-bottom: 1px;
    user-select: none;
  }
  .grid-row:not(.selected):hover {
    background-color: #f5f5f5;
  }
  .selected {
    background-color: #eee;
  }
  .div-cell {
    height: 25px;
    white-space: nowrap;
    overflow-x: hidden;
    display: flex;
    align-items: center;
  }
  .name {
    font-weight: 500;
    padding-left: 8px;
  }
  .bar {
    position: relative;
  }
  .value {
    padding-left: 10px;
    font-weight: 500;
  }
  .desc {
    text-overflow: ellipsis;
    padding-right: 5px;
  }
</style>
