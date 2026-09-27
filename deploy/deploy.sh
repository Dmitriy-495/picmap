#!/bin/bash
set -e

echo "🚀 Deploying PicMap..."

# 1. Собираем бинарник
echo "📦 Building backend..."
cd backend
go build -o picmap-backend .

# 2. Пушим исходники в GitHub
echo "📤 Pushing sources to GitHub..."
cd ..
git add .
git commit -m "deploy: $(date '+%Y-%m-%d %H:%M:%S')" || echo "Nothing to commit"
git push github main

# 3. Копируем бинарник на VPS
echo "📤 Copying binary to VPS..."
scp backend/picmap-backend pm:/var/www/picmap/backend/picmap-backend

# 4. Перезапускаем сервис на VPS
echo "🔄 Restarting service on VPS..."
ssh pm "sudo systemctl restart picmap-backend"

# 5. Проверяем
echo "✅ Checking health..."
sleep 2
curl -s http://picmap.ru/api/health || echo "⚠️  Health check failed"

echo ""
echo "✅ Deploy complete!"
