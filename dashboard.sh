#!/bin/bash

# --- 1. Setup ---
hour=$(date +%H)
greeting="Good Evening"
if [ "$hour" -lt 12 ]; then
    greeting="Good Morning"
elif [ "$hour" -lt 18 ]; then
    greeting="Good Afternoon"
fi

# --- 2. Banner ---
MAGENTA='\033[0;35m'
WHITE='\033[1;37m'
GREEN='\033[0;32m'
RED='\033[0;31m'
DARKGRAY='\033[1;30m'
YELLOW='\033[0;33m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

echo -e "${MAGENTA}"
echo "                                                   "
echo "     _____                              .___       "
echo "   _/ ____\____ ________ ____  ____   __| _/____   "
echo "   \   __\/ __ \___   // ___\/  _ \ / __ |/ __ \  "
echo "    |  | \  ___/ /    /\  \__(  <_> ) /_/ \  ___/  "
echo "    |__|  \___  >_____ \___  >____/\____ |\___  > "
echo "              \/      \/    \/           \/    \/  "
echo "                                                   "
echo -e "${NC}"

echo -e "  ${WHITE}${greeting}, ${USER}! Here are your top 3 TODO items, might want to focus on them:
${NC}"

# --- 3. Tasks ---
if command -v atlas.todo >/dev/null 2>&1; then
    # Capture output
    taskList=$(atlas.todo list desc 3 2>/dev/null)
    
    if [ -z "$taskList" ] || echo "$taskList" | grep -q "No pending tasks"; then
        echo -e "  ${GREEN}✨ Your board is clear! Ready for something new?${NC}"
    else
        echo "$taskList" | while IFS= read -r line; do
            if [[ "$line" == "[!"* ]]; then
                echo -e "  ${RED}${line}  🔥${NC}"
            elif [[ "$line" == "[."* ]]; then
                echo -e "  ${DARKGRAY}${line}${NC}"
            else
                echo -e "  ${YELLOW}${line}${NC}"
            fi
        done
    fi
else
    echo -e "  ${RED}[!] 'atlas.todo' command not found in PATH.${NC}"
fi

# --- 4. Footer ---
echo -e "
  ${DARKGRAY}--------------------------------------------------${NC}"
echo -e "  ${GRAY}Run 'atlas.todo' to manage tasks.${NC}
"