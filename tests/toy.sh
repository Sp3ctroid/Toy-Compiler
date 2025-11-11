#!/bin/bash

SOURCE_DIR="source"
LLVM_DIR="llvm_outputs"
EXEC_DIR="executables"

# Создаем папки если их нет
mkdir -p "$LLVM_DIR" "$EXEC_DIR"

# Компилируем все .txt файлы
for source_file in "$SOURCE_DIR"/*.txt; do
    if [ ! -f "$source_file" ]; then
        continue
    fi
    
    input_name=$(basename "$source_file" .txt)
    ll_file="$LLVM_DIR/$input_name.ll"
    exe_file="$EXEC_DIR/$input_name.exe"
    
    echo "Компиляция $input_name..."
    
    # Запуск compiler.exe
    ../compiler.exe -FN=$source_file -ON=$ll_file
    
    if [ -f "$ll_file" ]; then
        # Компиляция через Clang
        clang -o "$exe_file" "$ll_file"
        echo "Создан исполняемый файл: $exe_file"
    else
        echo "Ошибка: не удалось создать LLVM IR для $input_name"
    fi
done