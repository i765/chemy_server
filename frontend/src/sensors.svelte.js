import { config } from "./config.svelte";

class Sensors {
  // Глобальное реактивное состояние
  data = $state([]); // сырые данные с ацп (аналоговые и дискретные)
  analog = $state([]);
  isLoading = $state(true);
  error = $state(null);

  #timerId = null;
  #subscribers = 0;

  async fetchNow() {
    try {
      const response = await fetch("http://127.0.0.1:9765/frame");
      if (!response.ok) {
        throw new Error(`Ошибка HTTP: ${response.status}`);
      }
      this.data = await response.json();
      this.error = null;
    } catch (err) {
      this.error = err.message;
    } finally {
      this.isLoading = false;
      this.ConvertAnalog();
    }
  }

  startPolling(intervalMs = 1000) {
    // Запускаем таймер только если он еще не работает
    if (!this.timerId) {
      this.fetchNow(); // Первый запрос делаем сразу без задержки
      this.timerId = setInterval(() => {
        this.fetchNow();
      }, intervalMs);
    }

    // Возвращаем функцию остановки для $effect
    return () => {
      if (this.timerId) {
        clearInterval(this.timerId);
        this.timerId = null;
      }
    };
  }

  // сортировка и линейное преобразование
  ConvertAnalog() {
    for (let i = 0; i < config.data.analog.length; i++) {
      
      this.analog[i] = linearConvertAC(
        this.data.analogValues[i],
        config.data.analog[i].min,
        config.data.analog[i].max,
      );
      
    }

  }
  
}

// преобразование сигнала АЦП в диапазон шкалы измеряемого канала
function linearConvertAC(value = 0, min = 0, max = 1) {
  const round = roundFromRange(min, max);
  const dacMax = 1023;

  if (value <= 0) {
    return min;
  }

  if (value >= dacMax) {
    return max;
  }

  const converted = min + ((max - min) * value) / dacMax; // преобразуем

  return roundWithPrecision(converted, round);
}

// определение количества знаков после запятой
function roundFromRange(min = 0, max = 0) {
  let round = 2;
  let range = max - min;
  if (range >= 1000) {
    round = 0;
  } else if (range >= 100) {
    round = 1;
  }
  return round;
}

// корректное округление до заданного количества знаков
function roundWithPrecision(num, precision) {
  var multiplier = Math.pow(10, precision);
  return Math.round(num * multiplier) / multiplier;
}

export const sensors = new Sensors();
