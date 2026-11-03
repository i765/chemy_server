class VanillaCounter {
  constructor(containerSelector) {
    this.container = document.querySelector(containerSelector);
    this.count = 0;
    
    if (this.container) {
      this.render();
    }
  }

  increment() {
    this.count++;
    this.render(); // Перерисовываем при изменении состояния
  }

  render() {
    // Очищаем и наполняем контейнер заново
    this.container.innerHTML = `
      <div class="vanilla-counter">
        <button class="inc-btn">Добавить</button>
        <p>Кликнули: ${this.count} раз</p>
      </div>
    `;

    // Привязываем обработчик события
    this.container.querySelector('.inc-btn').addEventListener('click', () => this.increment());
  }
}

// Инициализация
const counter = new VanillaCounter('#counter-root');