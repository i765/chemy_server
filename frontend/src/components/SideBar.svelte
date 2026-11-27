<script>
  import { slide } from "svelte/transition";
  import { linear } from "svelte/easing";

  let items = [
    {
      id: "analog",
      label: "Аналоговые",
      open: false,
      children: [
        { id: "linear", label: "Линейные", href: "#section=analog&sub=linear" },
        { id: "calc", label: "Вычисляемые", href: "#section=analog&sub=calc" },
        { id: "deriv", label: "Производные", href: "#section=analog&sub=deriv" },
      ],
    },
    {
      id: "discrete",
      label: "Дискретные",
      open: false,
      children: [
        { id: "alarm", label: "Аварийные", href: "#section=discrete&sub=alarm" },
        { id: "g1", label: "Группа 1", href: "#section=discrete&sub=g1" },
        { id: "g2", label: "Группа 2", href: "#section=discrete&sub=g2" },
        { id: "g3", label: "Группа 3", href: "#section=discrete&sub=g3" },
        { id: "g4", label: "Группа 4", href: "#section=discrete&sub=g4" },
        { id: "g5", label: "Группа 5", href: "#section=discrete&sub=g5" },
      ],
    },
    {
      id: "moemos",
      label: "Мнемосхемы",
      open: false,
      children: [
        { id: "salt", label: "Солевые", href: "#salt" },
        { id: "wbp", label: "Переработка вод ЦБС", href: "#wbp" },
        { id: "gas", label: "Узел газоочистки", href: "#gas" },
        { id: "laundry", label: "Спецпрачечная", href: "#laundry" },
      ],
    },
  ];

  let activeId = $state(null); // Активный пункт
  let openStates = $state({}); // Состояние открытых разделов

  $effect(() => {
    activeId = localStorage.getItem("active") ?? "linear";
    openStates = JSON.parse(localStorage.getItem("openStates")) ?? {};
  });

  items.forEach((item) => {
    openStates[item.id] = item.open ?? false;
  });

  // Открыть / закрыть раздел
  function toggle(id) {
    openStates = { ...openStates, [id]: !openStates[id] };
    localStorage.setItem("openStates", JSON.stringify(openStates));
  }

  // Сделать пункт активным
  function activate(id) {
    activeId = id;
    localStorage.setItem("active", id);
  }
</script>

<aside class="sidebar">
  <nav class="sidebar-tree">
    <ul>
      {#each items as item (item.id)}
        <li class="tree-item">
          <button
            type="button"
            class="tree-parent"
            class:open={openStates[item.id]}
            aria-expanded={openStates[item.id]}
            aria-controls={item.id}
            onclick={() => toggle(item.id)}
          >
            {item.label}
            <span class="chevron"></span>
          </button>

          {#if openStates[item.id]}
            <ul
              id={item.id}
              class="tree-children"
              transition:slide={{ duration: 200, easing: linear }}
            >
              {#each item.children as child (child.id)}
                <li>
                  <a
                    href={child.href ?? "#"}
                    class="tree-child"
                    class:active={activeId === child.id}
                    onclick={() => activate(child.id)}
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
    --sidebar-bg: #fafafa;
    --sidebar-text: var(--main-text);
    --chevron-color: #999;

    --item-hover-text: #222;
    --item-hover-bg: #e5e5e5;
    --item-active-text: #111;
    --item-active-bg: #eee;

    background-color: var(--sidebar-bg);
    color: var(--sidebar-text);
    margin: 12px;
    padding: 12px 6px;
    display: flex;
    flex-direction: column;
    border-radius: 20px;
  }

  .sidebar-tree {
    padding: 0 12px;
  }

  .sidebar-tree ul {
    list-style: none;
  }

  .tree-parent {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 10px;
    background: none;
    border: none;
    text-align: left;
    cursor: pointer;
    font-family: inherit;
    font-weight: 500;
    transition:
      transform 200ms linear,
      opacity 200ms linear;
    opacity: 0.6;
  }

  .tree-parent:hover {
    opacity: 1;
  }

  .tree-parent .chevron {
    margin-left: auto;
    width: 8px;
    height: 8px;
    border-right: 2px solid var(--chevron-color);
    border-bottom: 2px solid var(--chevron-color);
    transform: rotate(-45deg);
    transition: transform 200ms linear;
  }

  .tree-parent.open .chevron {
    transform: rotate(45deg);
  }

  .tree-parent:hover .chevron {
    opacity: 1;
  }

  /* Дочерние пункты */
  .tree-children {
    padding-left: 16px;
  }

  .tree-child {
    display: flex;
    align-items: center;
    margin: 5px 0;
    padding: 8px 15px;
    text-decoration: none;
    border-radius: 8px;
    color: var(--sidebar-text);
    font-weight: 300;
  }

  .tree-child:hover {
    background-color: var(--item-hover-bg);
    color: var(--item-hover-text);
    transition: background-color 250ms ease-out;
  }

  .tree-child.active {
    background-color: var(--item-active-bg);
    color: var(--item-active-text);
    font-weight: 400;
    transition: background-color 250ms ease-out;
  }
</style>
