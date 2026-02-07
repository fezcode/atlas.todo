# --- 1. Setup ---
$hour = (Get-Date).Hour
$greeting = 'Good Evening'
if ($hour -lt 12) { $greeting = 'Good Morning' }
elseif ($hour -lt 18) { $greeting = 'Good Afternoon' }

# --- 2. Banner ---
Write-Host '                                                   ' -ForegroundColor Magenta
Write-Host '     _____                              .___       ' -ForegroundColor Magenta 
Write-Host '   _/ ____\____ ________ ____  ____   __| _/____   ' -ForegroundColor Magenta
Write-Host '   \   __\/ __ \\___   // ___\/  _ \ / __ |/ __ \  ' -ForegroundColor Magenta
Write-Host '    |  | \  ___/ /    /\  \__(  <_> ) /_/ \  ___/  ' -ForegroundColor Magenta
Write-Host '    |__|  \___  >_____ \\___  >____/\____ |\___  > ' -ForegroundColor Magenta
Write-Host '              \/      \/    \/           \/    \/  ' -ForegroundColor Magenta
Write-Host '                                                   ' -ForegroundColor Magenta

Write-Host "`n  $greeting, $($env:USERNAME)! Here are your top 3 TODO items, might want to focus on them:`n" -ForegroundColor White

# --- 3. Tasks ---
if (Get-Command atlas.todo -ErrorAction SilentlyContinue) {
    $taskList = & atlas.todo list desc 3

    if ($null -eq $taskList -or $taskList.Count -eq 0 -or ($taskList -join ' ') -match 'No pending tasks') {
        Write-Host '  ✨ Your board is clear! Ready for something new?' -ForegroundColor Green
    } else {
        foreach ($line in $taskList) {
            if ($line -like '[!]*') { 
                Write-Host "  $line" -ForegroundColor Red -NoNewline
                Write-Host '  🔥' -ForegroundColor Red
            }
            elseif ($line -like '[.]*') { 
                Write-Host "  $line" -ForegroundColor DarkGray 
            }
            else { 
                Write-Host "  $line" -ForegroundColor Yellow 
            }
        }
    }
} else {
    Write-Host "  [!] 'atlas.todo' command not found in PATH." -ForegroundColor Red
}

# --- 4. Footer ---
Write-Host "`n  --------------------------------------------------" -ForegroundColor DarkGray
Write-Host "  Run 't' or 'atlas.todo' to manage tasks." -ForegroundColor Gray
Write-Host ''

# --- 5. Alias ---
if (-not (Get-Alias t -ErrorAction SilentlyContinue)) {
    New-Alias -Name t -Value atlas.todo
}