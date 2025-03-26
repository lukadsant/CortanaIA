# **Instruções de Instalação e Uso**

## **Instalação**

1. **Mover a pasta de logs para o diretório correto**:
   - Mova a pasta `Logs` para o diretório `C://` do seu computador.
   
2. **Iniciar o script de execução**:
   - Execute o arquivo `start.vbs` para iniciar o processo.
   - O script irá rodar em segundo plano e, a cada **5 minutos**, o usuário verá um pop-up interativo.

3. **Ajustar o intervalo de tempo (opcional)**:
   - Se você quiser alterar o intervalo entre os pop-ups, modifique o valor do **cooldown** no arquivo `testex.bat` (em segundos).

---

## **Uso do Sistema**

### **1. Ver a mensagem atual**

Você pode verificar a mensagem atual através de um navegador, cURL ou Postman, utilizando o seguinte comando:
curl https://cortanaia.onrender.com






### **2. Salvar uma nova mensagem para a próxima vez que o programa for iniciado

curl -X POST https://cortanaia.onrender.com/post \
-H "Content-Type: application/json" \
-d '{
    "image_url": "https://static.wikia.nocookie.net/halo/images/c/c0/CortanaHalo5.png/revision/latest?cb=20211012223501&path-prefix=pt",
    "message": "oooooi tudo bem, essa é a implementação da nova versão da Cortana IA, posso realizar uma dúvida?"
}'




### **3. Ler o histórico de mensagens
Para visualizar o histórico de mensagens, use um navegador, cURL ou Postman:
curl https://cortanaia.onrender.com/logs

