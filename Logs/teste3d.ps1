Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

# Criar a janela do popup
$form = New-Object System.Windows.Forms.Form
$form.Text = "Cortana AI"
$form.Width = 500
$form.Height = 400
$form.StartPosition = "CenterScreen"
$form.FormBorderStyle = "FixedDialog"

# Definir ícone personalizado
$form.Icon = [System.Drawing.Icon]::ExtractAssociatedIcon("C:\Logs\assi.ico")

# Criar PictureBox para exibir imagem
$pictureBox = New-Object System.Windows.Forms.PictureBox
$pictureBox.Width = 100
$pictureBox.Height = 100
$pictureBox.Location = New-Object System.Drawing.Point(10, 10)
$pictureBox.SizeMode = "StretchImage"
$form.Controls.Add($pictureBox)

# Criar um rótulo para exibir a mensagem do assistente
$label = New-Object System.Windows.Forms.Label
$label.Text = "Carregando..."
$label.AutoSize = $true
$label.MaximumSize = New-Object System.Drawing.Size(370, 0)
$label.Location = New-Object System.Drawing.Point(120, 10)
$form.Controls.Add($label)

# Criar um campo de input (TextBox) para o usuário digitar a mensagem
$textBox = New-Object System.Windows.Forms.TextBox
$textBox.Location = New-Object System.Drawing.Point(20, 250)
$textBox.Width = 350
$form.Controls.Add($textBox)

# Criar um botão para enviar a mensagem do usuário
$sendButton = New-Object System.Windows.Forms.Button
$sendButton.Text = "Enviar"
$sendButton.Location = New-Object System.Drawing.Point(380, 250)
$form.Controls.Add($sendButton)

# Função para baixar imagem de um URL
function Download-Image {
    param ($url, $path)
    try {
        (New-Object System.Net.WebClient).DownloadFile($url, $path)
        return $path
    } catch {
        return $null
    }
}

# Função para quebrar texto automaticamente
function Format-Text {
    param ($text, $maxLength = 60)
    return ($text -split "(.{$maxLength})" | Where-Object {$_ -ne ""} ) -join "`n"
}

# Função para buscar a mensagem e imagem do assistente na API (GET)
function Get-AssistantMessage {
    $url = "https://cortanaia.onrender.com"
    try {
        $response = Invoke-RestMethod -Uri $url -Method Get
        return @{
            message = [System.Text.Encoding]::UTF8.GetString([System.Text.Encoding]::Default.GetBytes($response.message.ToString()))
            image   = $response.image_url
        }
    } catch {
        return @{ message = "Erro ao obter resposta do assistente."; image = "" }
    }
}

# Função para enviar o texto do usuário para a API (POST)
function Send-UserMessage {
    param ($userMessage)
    $url = "https://cortanaia.onrender.com/text"
    try {
        $body = @{ message = $userMessage } | ConvertTo-Json -Compress
        $headers = @{ "Content-Type" = "application/json" }
        $response = Invoke-RestMethod -Uri $url -Method Post -Body $body -Headers $headers
        return [System.Text.Encoding]::UTF8.GetString([System.Text.Encoding]::Default.GetBytes($response.message.ToString()))
    } catch {
        return "Erro ao enviar mensagem para o assistente."
    }
}

# Função para exibir mensagem na interface e ajustar o tamanho da janela
function Display-AssistantMessage {
    param ($message)
    $formattedMessage = Format-Text -text $message
    $label.Text = $formattedMessage
    $labelHeight = $formattedMessage.Split("`n").Count * 20
    $inputAreaHeight = 80
    $form.Height = [Math]::Max(400, 120 + $labelHeight + $inputAreaHeight)
}

# Função para exibir a animação "Pensando..." sem travar a UI
function Show-ThinkingAnimation {
    $thinkingStates = @("Pensando.", "Pensando..", "Pensando...")
    for ($i = 0; $i -lt 6; $i++) {
        $form.Invoke({ $label.Text = $thinkingStates[$i % $thinkingStates.Length] })
        Start-Sleep -Milliseconds 500
    }
}

# Buscar a mensagem inicial e imagem do assistente
$assistantData = Get-AssistantMessage
Display-AssistantMessage -message $assistantData.message

# Baixar e exibir imagem do assistente
if ($assistantData.image -ne "") {
    $localImagePath = "C:\Logs\assistant_image.jpg"
    $imagePath = Download-Image -url $assistantData.image -path $localImagePath
    if ($imagePath) {
        $pictureBox.Image = [System.Drawing.Image]::FromFile($imagePath)
    }
}

# Ação do botão "Enviar"
$sendButton.Add_Click({
    $userMessage = $textBox.Text
    if ($userMessage -ne "") {
        $formattedUserMessage = Format-Text -text "Você: $userMessage"
        Display-AssistantMessage -message $formattedUserMessage
        $textBox.Clear()

        # Mostrar animação "Pensando..." sem travar a interface
        Start-Job -ScriptBlock { Show-ThinkingAnimation }

        # Enviar mensagem para a API e receber resposta
        $assistantReply = Send-UserMessage -userMessage $userMessage
        Display-AssistantMessage -message "Assistente: $assistantReply"

        # Aguardar 5 segundos e fechar a janela automaticamente
        Start-Sleep -Seconds 5
        $form.Close()
    }
})

# Exibir o popup
$form.ShowDialog()
