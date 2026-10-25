#!/bin/bash
set -e

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}🚀 Ship PicMap (save + deploy)${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""

# ══════════════════════════════════════
# 1. GIT SAVE
# ══════════════════════════════════════
echo -e "${YELLOW}📋 Git status:${NC}"
git status --short
echo ""

HAS_CHANGES=false
if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$(git ls-files --others --exclude-standard)" ]; then
    HAS_CHANGES=true
fi

if [ "$HAS_CHANGES" = true ]; then
    echo -e "${YELLOW}📦 Добавляем изменения...${NC}"
    git add .
    echo ""

    echo -e "${YELLOW}📝 Формируем сообщение...${NC}"

    SCOPE="project"
    if git diff --cached --name-only | grep -q "frontend/"; then SCOPE="frontend"
    elif git diff --cached --name-only | grep -q "backend/"; then SCOPE="backend"
    elif git diff --cached --name-only | grep -q "docs/"; then SCOPE="docs"
    fi

    FILES=$(git diff --cached --name-only | head -10 | sed 's/^/  - /')
    FILE_COUNT=$(git diff --cached --name-only | wc -l | tr -d ' ')
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

    COMMIT_MSG="save(${SCOPE}): ${FILE_COUNT} file(s) @ ${TIMESTAMP}

${FILES}"

    if [ "$FILE_COUNT" -gt 10 ]; then
        COMMIT_MSG="${COMMIT_MSG}
  ... и ещё $((FILE_COUNT - 10))"
    fi

    echo -e "${YELLOW}💾 Коммитим...${NC}"
    git commit -m "$COMMIT_MSG"
    echo ""
else
    echo -e "${GREEN}✅ Изменений нет — пропускаем коммит${NC}"
    echo ""
fi

# ══════════════════════════════════════
# 2. GIT PUSH
# ══════════════════════════════════════
echo -e "${YELLOW}📤 Пушим на GitHub...${NC}"
git push github main
echo ""

# ══════════════════════════════════════
# 3. DEPLOY: BACKEND
# ══════════════════════════════════════
echo -e "${BLUE}────────────────────────────────────────${NC}"
echo -e "${YELLOW}📦 Building backend...${NC}"
cd backend
go build -o picmap-backend .
cd ..
echo ""

echo -e "${YELLOW}⏸  Stopping backend...${NC}"
ssh pm "sudo systemctl stop picmap-backend" || true
echo ""

echo -e "${YELLOW}📤 Copying backend to VPS...${NC}"
gzip -c backend/picmap-backend | ssh pm "gunzip > /var/www/picmap/backend/picmap-backend.new && chmod +x /var/www/picmap/backend/picmap-backend.new && mv /var/www/picmap/backend/picmap-backend.new /var/www/picmap/backend/picmap-backend"
echo ""

echo -e "${YELLOW}▶️  Starting backend...${NC}"
ssh pm "sudo systemctl start picmap-backend"
echo ""

# ══════════════════════════════════════
# 4. DEPLOY: FRONTEND
# ══════════════════════════════════════
echo -e "${BLUE}────────────────────────────────────────${NC}"
echo -e "${YELLOW}🎨 Building frontend (SSR)...${NC}"
cd frontend
pnpm run build
cd ..
echo ""

echo -e "${YELLOW}⏸  Stopping frontend...${NC}"
ssh pm "sudo systemctl stop picmap-frontend" || true
echo ""

echo -e "${YELLOW}📤 Copying frontend SSR to VPS...${NC}"
ssh pm "rm -rf /var/www/picmap/frontend-ssr/*"
rsync -avz --delete frontend/.output/ pm:/var/www/picmap/frontend-ssr/
echo ""

echo -e "${YELLOW}▶️  Starting frontend...${NC}"
ssh pm "sudo systemctl start picmap-frontend"
echo ""

# ══════════════════════════════════════
# 5. HEALTH CHECK
# ══════════════════════════════════════
echo -e "${YELLOW}⏳ Waiting...${NC}"
sleep 3
echo ""

echo -e "${YELLOW}✅ Checking backend...${NC}"
if curl -sf http://picmap.ru/api/health > /dev/null; then
    echo -e "${GREEN}✅ Backend OK${NC}"
else
    echo -e "${RED}❌ Backend failed${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}✅ Checking frontend...${NC}"
if curl -sf http://picmap.ru/ > /dev/null; then
    echo -e "${GREEN}✅ Frontend OK${NC}"
else
    echo -e "${RED}❌ Frontend failed${NC}"
    ssh pm "sudo journalctl -u picmap-frontend -n 30 --no-pager"
    exit 1
fi
echo ""

echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ SHIPPED!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo -e "Commit: $(git log -1 --format='%h %s')"
echo -e "Backend: http://picmap.ru/api/health"
echo -e "Frontend: http://picmap.ru/"
echo ""
