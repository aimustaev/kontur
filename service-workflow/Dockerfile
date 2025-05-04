FROM node:18

WORKDIR /app

# Установка зависимостей
COPY package*.json ./
RUN npm install

# Копирование исходного кода
COPY . .

EXPOSE 3000

# Используем entrypoint.sh для выбора режима запуска
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

# По умолчанию запускаем API сервер
CMD ["api"] 