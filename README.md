
# PicMap

Приложение для создания маршрутов путешествий и формирования воспоминаний.

## Стек

- **Backend:** Go + Gin + PostgreSQL 17
- **Frontend:** Nuxt 4 + Nuxt UI (планируется)
- **Карта:** Leaflet (планируется)
- **Миграции:** golang-migrate
- **Деплой:** GitHub + scp + systemd

## Структура
picmap/
├── docs/ # Документация
├── backend/ # Go-бэкенд
├── frontend/ # Nuxt-фронтенд (планируется)
├── deploy/ # Скрипты деплоя
└── README.md

text

## Разработка

### Backend

```bash
cd backend

# Запуск локально
go run main.go

# Сборка
go build -o picmap-backend .

# Миграции (на VPS)
migrate -path ./migrations -database "$DB_URL" up
Деплой
bash
./deploy/deploy.sh
Прод
VPS: 88.218.67.8 (Cloud.ru, Ubuntu 24.04)

Домен: picmap.ru

Backend: /var/www/picmap/backend/

Systemd: picmap-backend.service

Nginx: /etc/nginx/sites-available/picmap.ru

Git: GitHub Dmitriy-495/picmap

