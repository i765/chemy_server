class Config {
  items = $state([]);
  isLoading = $state(true);
  error = $state(null);

  async fetchNow() {
    try {
      const response = await fetch("http://127.0.0.1:9765/config");
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

}

export const config = new Config();