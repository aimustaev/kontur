# Service Kontur Frontend

Фронтенд сервис для системы тикетов Контур, построенный на React и Ant Design.

## Требования

- Node.js 20 или выше
- npm 10 или выше

## Установка

1. Установите зависимости:
```bash
npm install
```

2. Запустите сервер разработки:
```bash
npm run dev
```

Приложение будет доступно по адресу http://localhost:5173

## Сборка для продакшена

Для сборки приложения выполните:
```bash
npm run build
```

Собранные файлы будут находиться в директории `dist/`.

## Docker

Для сборки Docker образа:
```bash
docker build -t service-kontur-frontend .
```

Для запуска контейнера:
```bash
docker run -p 80:80 service-kontur-frontend
```

## Kubernetes

Для деплоя в Kubernetes используйте манифест из директории `k8s/frontend-deployment.yaml`:

```bash
kubectl apply -f k8s/frontend-deployment.yaml
```

## Структура проекта

- `src/` - исходный код приложения
- `public/` - статические файлы
- `k8s/` - Kubernetes манифесты
- `nginx.conf` - конфигурация Nginx
- `Dockerfile` - инструкции для сборки Docker образа

## Разработка

- `npm run dev` - запуск сервера разработки
- `npm run build` - сборка для продакшена
- `npm run lint` - проверка кода линтером
- `npm run preview` - предпросмотр собранного приложения 