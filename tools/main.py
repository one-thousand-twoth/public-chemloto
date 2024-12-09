import csv
import json
import os
import pathlib
import pandas as pd 

def process_csv_to_json(csv_file):
    """
    Обрабатывает CSV-файл и преобразует его в JSON-структуру.
    """
    with open(csv_file, 'r', encoding='utf-8') as f:
        reader = csv.reader(f)
        data = list(reader)
    # data = [row for row in data if row]
    # Переориентируем данные: первый столбец - ключи, последующие столбцы - значения
    keys = [row[0] for row in data]  # Первый столбец как ключи
    columns = zip(*[row[1:] for row in data if row.count("0") != len(row )- 1])  # Оставшиеся столбцы как данные

    # Создаем список словарей
    json_data = []
    for col in columns:
        # json_data.append({keys[i]: int(value) if value.isdigit() else 0 for i, value in enumerate(col)})
        col_values = [int(value) if value.isdigit() else 0 for value in col]
        if any(value != 0 for value in col_values):
            json_data.append({keys[i]: col_values[i] for i in range(len(keys))})
    # json_data = []
    # for col in range(1, len(data[0])):
    #     json_data.append({row[0]: int(row[col]) if row[col].isdigit() else 0 for row in data})
    return json_data

def process_all_csv_in_directory(directory):
    """
    Обрабатывает все CSV-файлы в указанной директории и сохраняет JSON-результаты.
    """
    result = {}
    osfiles = sorted(
        os.listdir(directory),
        key=lambda x: x[:3]
    )
    for folder in osfiles:
        result[folder] = {}
        files = sorted(
            os.listdir(os.path.join(directory, folder)),
            key=lambda x: x[:3])
        for file in files:
            try:
                result[folder][file.replace('.csv', '')[3:]] = process_csv_to_json(os.path.join(directory,folder, file))
            except Exception as e:
                print(f"Ошибка при обработке файла {file}: {e}")
    # for file in osfiles:
    #     if file.endswith('.csv'):
    #         file_path = os.path.join(directory, file)
    #         try:
    #             result[file.replace('.csv', '')[3:]] = process_csv_to_json(file_path)
    #         except Exception as e:
    #             print(f"Ошибка при обработке файла {file}: {e}")
    # Сохраняем результат в JSON-файл
    output_file = os.path.join(pathlib.Path(__file__).parent.__str__(), "output.json")
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, indent=4, ensure_ascii=False)
    print(f"Результаты сохранены в {output_file}")


# # Create JSON
# json_data = []
# for col in range(1, len(data[0])):
#     json_data.append({row[0]: int(row[col]) if row[col].isdigit() else 0 for row in data})

# # Print JSON
# formatted_json = json.dumps(json_data, indent=4)
# # print(formatted_json)

# with  open("struct.json", "w") as f:
#     f.write(formatted_json)

csv_directory = pathlib.Path(__file__).parent.joinpath("csv","polymers")
process_all_csv_in_directory(csv_directory.__str__())