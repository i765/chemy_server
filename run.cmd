@echo off

:: Быстрый запуск без вопросов (YES / NO)
set FASTSTART=NO

:: Последовательный порт
set SERIALPORT=COM1

:: Скорость чтения последовательного порта
set BAUD=19200

:: NONE, EVEN, ODD, MARK, SPACE
set PARITY=NONE 

:: Размер сглаживающего буфера, кадров
set BUFFERSIZE=16



.\chemyserver.exe