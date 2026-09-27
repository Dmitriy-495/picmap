#!/bin/bash
set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}🚀 Deploying PicMap${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""

# 1. Собираем бинарник
echo -e "${YELLOW}📦 Building backend...${NC}"
cd backend
go build -o picmap-backend .
cd ..

# 2. Пушим исходники в GitHub
echo ""
echo -e "${YELLOW}📤 Pushing sources to GitHub...${NC}"
git add .
if git diff --cached --quiet; then
    echo "  Nothing to commit"
else
    git commit -m "deploy: $(date '+%Y-%m-%d %H:%M:%S')"
fi
git push github main

# 3. Останавливаем сервис на VPS
echo ""
echo -e "${YELLOW}⏸  Stopping service on VPS...${NC}"
ssh pm "sudo systemctl stop picmap-backend" || true

# 4. Копируем бинарник на VPS
echo ""
echo -e "${YELLOW}📤 Copying binary to VPS...${NC}"
scp backend/picmap-backend pm:/var/www/picmap/backend/picmap-backend

# 5. Запускаем сервис
echo ""
echo -e "${YELLOW}▶️  Starting service on VPS...${NC}"
ssh pm "sudo systemctl start picmap-backend"

# 6. Ждём и проверяем
echo ""
echo -e "${YELLOW}⏳ Waiting for service to start...${NC}"
sleep 3

echo ""
echo -e "${YELLOW}✅ Checking health...${NC}"
if curl -sf http://picmap.ru/api/health; then
    echo ""
    echo -e "${GREEN}════════════════════════════════════════${NC}"
    echo -e "${GREEN}✅ Deploy complete!${NC}"
    echo -e "${GREEN}════════════════════════════════════════${NC}"
else
    echo ""
    echo -e "${RED}════════════════════════════════════════${NC}"
    echo -e "${RED}❌ Health check failed!${NC}"
    echo -e "${RED}════════════════════════════════════════${NC}"
    echo ""
    echo "Check logs on VPS:"
    echo "  ssh pm 'sudo journalctl -u picmap-backend -n 50'"
    exit 1
fi
