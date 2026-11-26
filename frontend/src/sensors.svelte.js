import { config } from "./config.svelte";

class Sensors {
  // Глобальное реактивное состояние
  data = $state([]); // сырые данные с УСО (аналоговые и дискретные)
  analog = $state([]); // Преобразованные аналоговые данные
  isReady = $state(false); // Данные готовы
  error = $state(null);

  #timerId = null;
  #subscribers = 0;

  async fetchNow() {
    try {
      const response = await fetch("http://127.0.0.1:9765/frame");
      if (response.ok) {
        this.data = await response.json();
        this.error = null
      }

    } catch (err) {
      this.error = err;
      this.isReady = false;

    } finally {
      if (this.error == null) {       
        if (config.isReady) {
          this.ConvertAnalog(); // Если конфигурация готова, работаем с данными
          this.isReady = true;
        } else {
          config.fetchNow(); // Если нет, делаем запрос конфигурации
        }
      }
      
    }
  }

  startPolling(intervalMs = 1000) {
    // Запускаем интервал только если он еще не работает
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
    const channelOffset = 1; // Смещение индексов каналов
    for (let i = 0; i < config.data.analog.length; i++) {
      const address = config.data.analog[i].address + channelOffset; // Адресация
      this.analog[i] = linearConvertAC(
        this.data.analogValues[address], // Адресация
        config.data.analog[i].min,
        config.data.analog[i].max,
        config.data.analog[i].factor,
        config.data.analog[i].offset,
      );
    }
  }
}

// преобразование сигнала АЦП в диапазон шкалы измеряемого канала
function linearConvertAC(value = 0, min = 0, max = 1, factor = 1, offset = 0) {
  const dacMax = 1023;
  
  // 1. Зажимаем сырое значение АЦП в допустимых пределах
  const clampedValue = Math.min(Math.max(value, 0), dacMax);

  // 2. Считаем по единой формуле (offset и factor сработают и на краях)
  const converted = (min + ((max - min) * clampedValue) / dacMax + offset) * factor;

  // 3. Точность лучше считать с учетом масштаба (factor)
  const round = roundFromRange(min * factor, max * factor);

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
