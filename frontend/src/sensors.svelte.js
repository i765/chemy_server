class Sensors {
  // Глобальное реактивное состояние
  items = $state([]);
  isLoading = $state(true);
  error = $state(null);

  #timerId = null;
  #subscribers = 0;

  async fetchNow() {
    try {
      const response = await fetch('http://127.0.0.1:9765/frame');
      if (!response.ok) {
        throw new Error(`Ошибка HTTP: ${response.status}`);
      }
      this.items = await response.json();
      this.error = null;
    } catch (err) {
      this.error = err.message;
    } finally {
      this.isLoading = false;
    }
  }

  // Запуск опроса (с подсчетом активных компонентов)
  startPolling(intervalMs = 1000) {
    this.#subscribers++;
    
    // Запускаем таймер только если он еще не работает
    if (!this.#timerId) {
      this.fetchNow(); // Первый запрос делаем сразу без задержки
      this.#timerId = setInterval(() => {
        this.fetchNow();
      }, intervalMs);
    }

    // Возвращаем функцию остановки для $effect
    return () => {
      this.#subscribers--;
      if (this.#subscribers <= 0 && this.#timerId) {
        clearInterval(this.#timerId);
        this.#timerId = null;
        this.#subscribers = 0;
      }
    };
  }
}

export const sensors = new Sensors();
