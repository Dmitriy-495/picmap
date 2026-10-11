#!/bin/bash
set -e

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}🚀 Deploying PicMap${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""

# ══════════════════════════════════════
# 1. BACKEND
# ══════════════════════════════════════
echo -e "${YELLOW}📦 Building backend...${NC}"
cd backend
go build -o picmap-backend .
cd ..

# ══════════════════════════════════════
# 2. FRONTEND
# ══════════════════════════════════════
echo ""
echo -e "${YELLOW}🎨 Building frontend...${NC}"
cd frontend
pnpm run build
cd ..

# ══════════════════════════════════════
# 3. GIT PUSH
# ══════════════════════════════════════
echo ""
echo -e "${YELLOW}📤 Pushing sources to GitHub...${NC}"
git add .
if git diff --cached --quiet; then
    echo "  Nothing to commit"
else
    git commit -m "deploy: $(date '+%Y-%m-%d %H:%M:%S')"
fi
git push github main

# ══════════════════════════════════════
# 4. BACKEND → VPS
# ══════════════════════════════════════
echo ""
echo -e "${YELLOW}⏸  Stopping backend service...${NC}"
ssh pm "sudo systemctl stop picmap-backend" || true

echo ""
echo -e "${YELLOW}📤 Copying backend to VPS...${NC}"
gzip -c backend/picmap-backend | ssh pm "gunzip > /var/www/picmap/backend/picmap-backend.new && chmod +x /var/www/picmap/backend/picmap-backend.new && mv /var/www/picmap/backend/picmap-backend.new /var/www/picmap/backend/picmap-backend"

echo ""
echo -e "${YELLOW}▶️  Starting backend...${NC}"
ssh pm "sudo systemctl start picmap-backend"

# ══════════════════════════════════════
# 5. FRONTEND → VPS
# ══════════════════════════════════════
echo ""
echo -e "${YELLOW}📤 Copying frontend to VPS...${NC}"
ssh pm "mkdir -p /var/www/picmap/frontend"
rsync -avz --delete frontend/.output/public/ pm:/var/www/picmap/frontend/

# ══════════════════════════════════════
# 6. HEALTH CHECK
# ══════════════════════════════════════
echo ""
echo -e "${YELLOW}⏳ Waiting for services...${NC}"
sleep 3

echo ""
echo -e "${YELLOW}✅ Checking backend health...${NC}"
if curl -sf http://picmap.ru/api/health > /dev/null; then
    echo -e "${GREEN}✅ Backend OK${NC}"
else
    echo -e "${RED}❌ Backend failed${NC}"
    ssh pm "sudo journalctl -u picmap-backend -n 20"
    exit 1
fi

echo ""
echo -e "${YELLOW}✅ Checking frontend...${NC}"
if curl -sf http://picmap.ru/ > /dev/null; then
    echo -e "${GREEN}✅ Frontend OK${NC}"
else
    echo -e "${RED}❌ Frontend failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Deploy complete!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
