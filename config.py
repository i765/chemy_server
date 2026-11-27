import json
import shutil
from pathlib import Path

config_path = Path("config.json")
backup_path = Path("config.json.bak")

# Создаём резервную копию перед изменением
shutil.copy2(config_path, backup_path)

with config_path.open("r", encoding="utf-8") as file:
    data = json.load(file)

for item in data.get("analog", []):
    max_val = item.get("max", 0)
    min_val = item.get("min", 0)
    span = max_val - min_val

    item["factor"] = 1
    item["offset"] = 0
    item["red_h"] = round(min_val + span * 0.95, 2)
    item["yel_h"] = round(min_val + span * 0.85, 2)
    item["yel_l"] = round(min_val + span * 0.025, 2)
    item["red_l"] = round(min_val + span * 0.01, 2)

with config_path.open("w", encoding="utf-8") as file:
    json.dump(data, file, indent=2, ensure_ascii=False)

print(f"Готово! Обновлено датчиков: {len(data.get('analog', []))}")
