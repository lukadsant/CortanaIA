@echo off
:loop
rem Executa o script VBS em segundo plano, sem mostrar a janela de terminal
start "" /B wscript C:\Logs\runteste.vbs
rem Espera 5 minutos (60 segundos)
ping -n 301 127.0.0.1 > nul
goto loop
