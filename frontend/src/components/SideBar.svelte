<script>
  import { createEventDispatcher } from 'svelte';

  // Данные меню. Можно передать снаружи через `items`.
  export let items = [
    {
      id: 'analog',
      label: 'Аналоговые',
      open: true,
      children: [
        { id: 'linear', label: 'Линейные', href: '#linear' },
        { id: 'calc',   label: 'Вычисляемые', href: '#calc' },
        { id: 'deriv',  label: 'Производные', href: '#deriv' },
      ],
    },
    {
      id: 'discrete',
      label: 'Дискретные',
      open: false,
      children: [
        { id: 'alarm', label: 'Аварийные', href: '#alarm' },
        { id: 'g1', label: 'Группа 1', href: '#g1' },
        { id: 'g2', label: 'Группа 2', href: '#g2' },
        { id: 'g3', label: 'Группа 3', href: '#g3' },
        { id: 'g4', label: 'Группа 4', href: '#g4' },
        { id: 'g5', label: 'Группа 5', href: '#g5' },
      ],
    },
    {
      id: 'moemos',
      label: 'Мнемосхемы',
      open: false,
      children: [
        { id: 'salt', label: 'Солевые', href: '#salt' },
        { id: 'wbp', label: 'Переработка вод ЦБС', href: '#wbp' },
        { id: 'gas', label: 'Узел газоочистки', href: '#gas' },
        { id: 'laundry', label: 'Спецпрачечная', href: '#laundry' },
      ],
    },
  ];

  // Активный пункт (можно задать снаружи)
  export let activeId = null;

  // Диспетчер — через него состояние экспортируется наружу
  const dispatch = createEventDispatcher();

  // Состояние открытых разделов (инициализируем из данных)
  let openStates = {};
  items.forEach((item) => {
    openStates[item.id] = item.open ?? false;
  });

  // Открыть / закрыть раздел
  function toggle(id) {
    openStates = { ...openStates, [id]: !openStates[id] };
    dispatch('toggle', { id, open: openStates[id] });
    dispatch('change', { openStates, activeId });
  }

  // Сделать пункт активным
  function activate(id) {
    activeId = id;
    dispatch('activate', { id });
    dispatch('change', { openStates, activeId });
  }
</script>

<aside class="sidebar">
  <nav class="sidebar-tree" aria-label="Навигация по разделам">
    <ul>
      {#each items as item (item.id)}
        <li class="tree-item">
          <button
            type="button"
            class="tree-parent"
            class:open={openStates[item.id]}
            aria-expanded={openStates[item.id]}
            aria-controls={item.id}
            on:click={() => toggle(item.id)}
          >
            {item.label}
            <span class="chevron" aria-hidden="true"></span>
          </button>

          {#if openStates[item.id]}
            <ul id={item.id} class="tree-children">
              {#each item.children as child (child.id)}
                <li>
                  <a
                    href={child.href ?? '#'}
                    class="tree-child"
                    class:active={activeId === child.id}
                    on:click={() => activate(child.id)}
                  >
                    {child.label}
                  </a>
                </li>
              {/each}
            </ul>
          {/if}
        </li>
      {/each}
    </ul>
  </nav>
</aside>

<style>
  .sidebar {
    background-color: var(--sidebar-bg);
    color: var(--sidebar-text);
    margin: 10px;
    padding: 10px;
    display: flex;
    flex-direction: column;
    border-radius: 20px;
  }

  .sidebar-tree {
    max-width: 300px;
    padding: 0 10px;
  }

  .sidebar-tree ul {
    list-style: none;
  }
  /* Родительские пункты — теперь это button (доступность + клик) */
  .tree-parent {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 10px 14px;
    color: var(--sidebar-text);
    background: none;
    border: none;
    font-size: 15px;
    font-weight: 500;
    text-align: left;
    cursor: pointer;
    font-family: inherit;
  }

  /* Стрелка-шеврон — отдельный элемент, анимируем transform */
  .tree-parent .chevron {
    margin-left: auto;
    width: 8px;
    height: 8px;
    border-right: 2px solid var(--line-color);
    border-bottom: 2px solid var(--line-color);
    transition: transform .2s ease;
    transform: rotate(45deg);
  }

  .tree-parent.open .chevron {
    transform: rotate(-135deg);
  }

  /* Дочерние пункты */
  .tree-children {
    padding-left: 16px;
    border-left: 1px dashed var(--line-color);
  }

  .tree-child {
    display: flex;
    align-items: center;
    margin: 5px 0;
    padding: 6px 13px;
    font-size: 14px;
    color: #333333;
    font-weight: 400;
    text-decoration: none;
    border-radius: 5px;
  }

  .tree-child:hover {
    background-color: var(--line-color);
    color: #fff;
  }

  .tree-child.active {
    background-color: var(--line-color);
    color: #fff;
  }

</style>
