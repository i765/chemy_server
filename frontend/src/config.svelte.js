class Config {
  data = $state([]);
  isReady = $state(false);
  error = $state(null);
  
  async fetchNow() {
    try {
      const response = await fetch("http://127.0.0.1:9765/config");
      if (response.ok) {
        this.data = await response.json();
        this.error = null
      }
      
    } catch (err) {
      this.error = err
    } finally {
      if (this.error === null) {
        this.isReady = true;
        console.log('Конфигурация загружена.')
      }
    }
  }
}

export const config = new Config();
