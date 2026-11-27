import { config } from "./config.svelte";

class Sensors {
  // Глобальное реактивное состояние
  data = $state([]); // сырые данные с УСО (аналоговые и дискретные)
  analog = $state([]); // Преобразованные аналоговые данные
  warnings = $state([]);

  isReady = $state(false); // Данные готовы
  error = $state(null);

  #timerId = null;
  #subscribers = 0;

  async fetchNow() {
    try {
      const response = await fetch("http://127.0.0.1:9765/frame");
      if (response.ok) {
        this.data = await response.json();
        this.error = null;
      }
    } catch (err) {
      this.error = err;
      this.isReady = false;
    } finally {
      if (this.error == null) {
        if (config.isReady) {
          this.processAnalog(); // Если конфигурация готова, отправляем данные на обработку
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

  // обработка аналога
  processAnalog() {
    const channelOffset = 1; // Смещение индексов каналов
    for (let i = 0; i < config.data.analog.length; i++) {
      const address = config.data.analog[i].address + channelOffset; // Адресация
      const value = linearConvertAC(
        this.data.analogValues[address], // Адресация
        config.data.analog[i].min,
        config.data.analog[i].max,
        config.data.analog[i].factor,
        config.data.analog[i].offset,
      );

      this.warnings[i] = CheckStatus(
        value,
        config.data.analog[i].red_h ?? null,
        config.data.analog[i].yel_h ?? null,
        config.data.analog[i].yel_l ?? null,
        config.data.analog[i].red_l ?? null,
      );

      this.analog[i] = value
    }
  }
}

function CheckStatus(value = 0, red_h = null, yel_h = null, yel_l = null, red_l = null) {
    // 1. Сначала проверяем критические пороги (Красный - 2), так как у аварии высший приоритет
    if ((red_h !== null && value >= red_h) || (red_l !== null && value <= red_l)) {
        return 2;
    }

    // 2. Затем проверяем предупредительные пороги (Желтый - 1)
    if ((yel_h !== null && value >= yel_h) || (yel_l !== null && value <= yel_l)) {
        return 1;
    }

    // 3. Если ни один из порогов не нарушен, возвращаем норму (0)
    return 0;
}

// преобразование сигнала АЦП в диапазон шкалы измеряемого канала
function linearConvertAC(value = 0, min = 0, max = 1, factor = 1, offset = 0) {
  const dacMax = 1023;
  // Зажимаем сырое значение АЦП в допустимых пределах
  const clampedValue = Math.min(Math.max(value, 0), dacMax);
  // Считаем по единой формуле (offset и factor сработают и на краях)
  const converted =
    (min + ((max - min) * clampedValue) / dacMax + offset) * factor;
  // Точность считаем с учетом масштаба (factor)
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
