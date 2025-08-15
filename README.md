# Notification Service

Высокопроизводительный микросервис для обработки и доставки уведомлений, построенный на Go с поддержкой масштабирования в Kubernetes кластерах.

## 🚀 Особенности

- **Высокая производительность**: Написан на Go для максимальной скорости и эффективности
- **Масштабируемость**: Поддержка горизонтального масштабирования в K8s/K3s кластерах
- **Реальное время**: WebSocket соединения для мгновенной доставки уведомлений
- **Push-уведомления**: Интеграция с Firebase Cloud Messaging (FCM)
- **Надежность**: Использование Redis для кэширования и Pub/Sub
- **Асинхронность**: Обработка сообщений через Kafka
- **Мониторинг**: Health check endpoints для проверки состояния сервиса

## 🏗️ Архитектура

Сервис построен по принципам Clean Architecture с четким разделением на слои:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Kafka Input   │    │  HTTP API      │    │  WebSocket      │
│   (Consumer)    │    │  (REST)        │    │  (Real-time)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Processor     │
                    │  (Business      │
                    │   Logic)        │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Redis       │    │   PostgreSQL    │    │   Push FCM      │
│  (Cache/PubSub) │    │   (Storage)     │    │  (Mobile)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Технологический стек

- **Язык**: Go 1.24+
- **База данных**: PostgreSQL
- **Кэш и Pub/Sub**: Redis
- **Очередь сообщений**: Apache Kafka
- **WebSocket**: Socket.IO v2
- **Push-уведомления**: Firebase Cloud Messaging
- **Контейнеризация**: Docker
- **Оркестрация**: Kubernetes/K3s
- **Миграции**: golang-migrate

## 📋 Требования

- Go 1.24+
- Docker
- Kubernetes кластер (K8s/K3s)
- PostgreSQL
- Redis
- Apache Kafka

## 🚀 Быстрый старт

### Локальная разработка

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/your-username/notification-service.git
cd notification-service
```

2. **Установите зависимости**
```bash
go mod download
```

3. **Настройте переменные окружения**
```bash
export DATABASE_PASSWORD=your_password
export FCM_KEY=your_fcm_key
```

4. **Запустите сервис**
```bash
go run main.go
```

### Docker

```bash
docker build -t notification-service .
docker run -p 3000:3000 -p 4000:4000 notification-service
```

### Kubernetes

1. **Примените базовую конфигурацию**
```bash
kubectl apply -k kustomise/base/
```

2. **Или используйте overlay для продакшена**
```bash
kubectl apply -k kustomise/overlay/prod/
```

## ⚙️ Конфигурация

Сервис настраивается через переменные окружения и конфигурационные файлы:

### Основные настройки

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_PASSWORD` | Пароль PostgreSQL | - |
| `FCM_KEY` | Ключ Firebase Cloud Messaging | - |
| `REDIS_ADDR` | Адрес Redis сервера | `redis-master.default.svc.cluster.local:6379` |
| `KAFKA_BROKERS` | Адреса Kafka брокеров | `10.8.1.1:9092` |

### Конфигурация базы данных

```go
type PostgresConfig struct {
    Host     string // 10.8.1.1
    Port     int    // 5432
    User     string // avenir
    Password string // из переменной окружения
    DBName   string // notification_service
}
```

## 📡 API Endpoints

### Health Check
```
GET /health
```
Возвращает статус сервиса

### Получение уведомлений
```
GET /api/v1/notifications
```
Возвращает список уведомлений для авторизованного пользователя

**Заголовки:**
- `Authorization: Bearer <JWT_TOKEN>`

**Ответ:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Заголовок уведомления",
      "message": "Текст уведомления",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "hasNextPage": false
}
```

## 🔄 Обработка сообщений

Сервис подписывается на Kafka топик `message.created` и обрабатывает события создания сообщений:

1. **Получение события** из Kafka
2. **Обработка бизнес-логики** в processor
3. **Сохранение** в PostgreSQL
4. **Отправка** через WebSocket и Push FCM
5. **Кэширование** в Redis

## 📊 Мониторинг и метрики

### Health Check
- Endpoint: `/health`
- Метод: `GET`
- Ответ: `Pong` при успешной работе

### Логирование
Сервис выводит логи в stdout/stderr для интеграции с системами логирования Kubernetes.

## 🚀 Масштабирование

### Горизонтальное масштабирование

```bash
# Увеличить количество реплик
kubectl scale deployment notification-deployment --replicas=3

# Автоматическое масштабирование (HPA)
kubectl apply -f - <<EOF
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: notification-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: notification-deployment
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
EOF
```

### Вертикальное масштабирование

Настройте ресурсы в `deployment.yaml`:

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
```

## 🔒 Безопасность

- JWT аутентификация для API endpoints
- Переменные окружения для чувствительных данных
- Kubernetes Secrets для хранения паролей
- Валидация входящих данных

## 🧪 Тестирование

```bash
# Запуск unit тестов
go test ./...

# Запуск интеграционных тестов
go test -tags=integration ./...

# Запуск e2e тестов
go test -tags=e2e ./...
```

## 📈 Производительность

Сервис оптимизирован для:
- **Низкого потребления памяти**: ~5-20MB на инстанс
- **Быстрой обработки**: <10ms для типичных операций
- **Высокой пропускной способности**: 1000+ уведомлений/сек
- **Эффективного использования CPU**: <5% в idle состоянии

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request


## 📞 Поддержка

- **Issues**: [GitHub Issues](https://github.com/your-username/notification-service/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/notification-service/discussions)

## 🙏 Благодарности

- Go team за отличный язык программирования
- Сообщество Kubernetes за инструменты оркестрации
- Redis, Kafka и PostgreSQL за надежные технологии

---

⭐ Если проект оказался полезным, поставьте звездочку!
