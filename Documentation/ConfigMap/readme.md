# ConfigMap в Kubernetes

Офицальная документация *ConfigMap*: [https://kubernetes.io/docs/concepts/configuration/configmap/](https://kubernetes.io/docs/concepts/configuration/configmap/)

---

#### **1. Что такое ConfigMap?**

**ConfigMap** — это объект Kubernetes для хранения **конфигурационных данных** в формате "ключ-значение". Основные особенности:

* Отделяет конфигурацию от образов контейнеров
* Позволяет изменять конфигурацию без пересборки образов
* Подходит для хранения:

  * Переменных среды (environment variables)
  * Параметров командной строки
  * Конфигурационных файлов

**Примеры использования**:

* Настройки подключения к БД
* Флаги функциональности (feature flags)
* Параметры логирования

---

#### **2. Создание ConfigMap**

#### **Из командной строки**

```bash
# Из литералов
kubectl create configmap app-config \
  --from-literal=LOG_LEVEL=debug \
  --from-literal=MAX_CONNECTIONS=50

# Из файла
kubectl create configmap nginx-config \
  --from-file=nginx.conf
```

#### **Через YAML-манифест**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: game-config
data:
  # Простые значения
  PLAYER_INITIAL_LIVES: "3"
  UI_PROPERTIES_FILE_NAME: "user-interface.properties"
  
  # Многострочные конфигурации
  game.properties: |
    enemy.types=aliens,monsters
    player.maximum-lives=5    
```

---

#### **3. Способы использования ConfigMap**

#### **Как переменные среды**

```yaml
containers:
- name: app
  image: my-app
  envFrom:
  - configMapRef:
      name: app-config
```

#### **Как отдельные переменные**

```yaml
env:
- name: LOG_LEVEL
  valueFrom:
    configMapKeyRef:
      name: app-config
      key: LOG_LEVEL
```

#### **Как файлы в Volume**

```yaml
volumes:
- name: config-volume
  configMap:
    name: game-config

volumeMounts:
- name: config-volume
  mountPath: /etc/config
```

---

#### **4. Обновление ConfigMap**

1. Внесите изменения:

    ```bash
    kubectl edit configmap app-config
    ```
2. Для Pod'ов:

    * **Переменные среды**: Требуется перезапуск Pod'а
    * **Смонтированные как Volume**: Обновляются автоматически (≈15-30 сек)

**Важно**: Для немедленного обновления можно:

* Удалить Pod для пересоздания
* Использовать сторонние решения (например, Reloader)

---

#### **5. Ограничения**

|Ограничение|Решение|
| ----------------------------------------------------| -------------------------------------------------------------|
|Макс. размер 1MiB|Разбивайте на несколько ConfigMap'ов|
|Нет версионирования|Используйте вместе с Helm|
|Нет встроенного шифрования|Для секретов используйте Secret|

---

#### **6. Best Practices**

1. **Именование**: Включайте версию/назначение (`app-config-v1`​)
2. **Структура**: Группируйте связанные настройки
3. **Чувствительные данные**: Всегда используйте **Secret** вместо ConfigMap
4. **Мониторинг**: Отслеживайте изменения через аудит
5. **Документация**: Комментируйте параметры в манифестах

---

#### **7. Пример: Полноценное приложение с ConfigMap**

```yaml
# Конфигурация
apiVersion: v1
kind: ConfigMap
metadata:
  name: web-app-config
data:
  APP_COLOR: blue
  APP_ENV: prod

# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  template:
    spec:
      containers:
      - name: app
        image: nginx
        envFrom:
        - configMapRef:
            name: web-app-config
```

---
