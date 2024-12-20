# Реализация команды `kill`

## Цель
Создать утилиту на Go, которая отправляет сигнал указанному процессу (аналог команды `kill` в Linux).

## Описание

1. **Приём аргументов:**
    - `-pid`: обязательный аргумент для указания PID процесса, которому нужно отправить сигнал.
    - `-signal`: опциональный аргумент для указания номера сигнала (например, `9` для SIGKILL). По умолчанию используется сигнал `SIGTERM`.

   **Пример запуска:**
   ```bash
   go-kill -pid 1234  
   # Отправит SIGTERM процессу с PID 1234  

   go-kill -pid 1234 -signal 9  
   # Отправит SIGKILL процессу с PID 1234
    ```
   
2. **Отправка сигнала:**
- Использовать пакеты os и syscall для отправки сигнала процессу.
- Обработать ошибки:
  - Процесс с указанным PID не существует.
  - Сигнал не поддерживается.

3. **Обработка флагов:**
        -pid: обязательный флаг для указания PID.
        -signal: номер сигнала (необязательный).

### Дополнительно (необязательно)
1. Поддержка текстовых сигналов:
 - Позволить указывать сигналы по имени (SIGKILL, SIGTERM) вместо числовых значений.
2. Проверка существования процесса перед отправкой сигнала.

### Подсказки
 - Используйте os.FindProcess для получения процесса по PID.
 - Для отправки сигналов пригодятся пакеты syscall и os/signal.

### iРезультат
Рабочая утилита, которая отправляет сигналы процессам, аналогичная команде kill в Linux.