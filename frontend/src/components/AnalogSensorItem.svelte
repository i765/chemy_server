<script>
  let { item, sensors } = $props();
</script>

<div class="div-cell id">{item.name}</div>
<div class="div-cell bar">asd</div>
<div class="div-cell">{sensors.analog}</div>
<div class="div-cell">{item.desc}</div>

<style>
  .div-cell {
    border-bottom: 1px solid #ccc;
  }
</style>
