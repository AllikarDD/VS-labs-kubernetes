# Аннотации (Annotations) в Kubernetes

Офицальная документация *Annotations*: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)

---

### **1. Что такое аннотации?**

**Аннотации** — это пары **ключ-значение**, которые используются для хранения произвольной метаинформации об объектах Kubernetes (например, Pod'ах, Services, Deployments).

В отличие от **меток (labels)** , аннотации:

* Не предназначены для фильтрации или группировки объектов.
* Могут содержать большие объёмы данных (например, описание, ссылки, конфигурации).
* Используются для внутренней логики инструментов (например, CI/CD, мониторинга).

---

### **2. Для чего используются аннотации?**

* **Хранение вспомогательной информации**:

  * Контактные данные владельца объекта.
  * Ссылки на документацию или тикеты (например, `issue-tracker/url: "https://example.com/ticket-123"`​).
* **Управление поведением инструментов**:

  * Указание для Ingress-контроллера, как обрабатывать трафик (`nginx.ingress.kubernetes.io/rewrite-target: /`​).
  * Настройка балансировки нагрузки (`service.beta.kubernetes.io/aws-load-balancer-type: external`​).
* **Логирование и аудит**:

  * Дата создания объекта (`created-by: "ci-pipeline-2023"`​).
  * Причина последнего обновления (`update-reason: "security-patch"`​).

---

### **3. Синтаксис аннотаций**

* **Ключи**:

  * Могут состоять из префикса (опционально) и имени, разделённых `/`​ (например, `company.com/log-format`​).
  * Префикс должен быть действительным DNS-поддоменом (как в метках).
  * Имя — до 63 символов (`[a-z0-9A-Z]`​, `-`​, `_`​, `.`​).
* **Значения**:

  * Любые строковые данные (включая JSON, XML, многострочный текст).
  * Нет жёстких ограничений на длину (но лучше избегать гигантских аннотаций).

**Пример аннотаций в Pod:**

```yaml
metadata:
  annotations:
    owner: "team-devops@example.com"
    commit-hash: "a1b2c3d"
    kubernetes.io/description: "This pod runs the main API service."
```

---

### **4. Разница между аннотациями и метками**

|Характеристика|Метки (Labels)|Аннотации (Annotations)|
| ------------------------------| -------------------------------------------------------| -------------------------------------------------------------------------------|
|**Назначение**|Группировка и выбор объектов|Хранение метаданных|
|**Использование**|Селекторы (`kubectl -l`​), Services|Документирование, настройка инструментов|
|**Ограничения**|Строгие правила именования|Более гибкие|
|**Пример**|​`env: production`​, `app: frontend`​|​`deployment.kubernetes.io/revision: "2"`​|

---

### **5. Практические примеры**

#### **Пример 1: Аннотации для Ingress-контроллера**

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: my-service
            port:
              number: 80
```

#### **Пример 2: Аннотации для мониторинга (Prometheus)**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  annotations:
    prometheus.io/scrape: "true"     # Разрешить сбор метрик
    prometheus.io/port: "8080"       # Порт для метрик
    prometheus.io/path: "/metrics"   # URL-путь
spec:
  containers:
  - name: app
    image: my-app:latest
```

---

### **6. Работа с аннотациями через** **​`kubectl`​**​

* **Добавить аннотацию:**

  ```sh
  kubectl annotate pods my-pod owner="team-devops@example.com"
  ```
* **Обновить аннотацию:**

  ```sh
  kubectl annotate pods my-pod owner="new-owner@example.com" --overwrite
  ```
* **Удалить аннотацию:**

  ```sh
  kubectl annotate pods my-pod owner-
  ```
* **Просмотреть аннотации:**

  ```sh
  kubectl get pod my-pod -o jsonpath='{.metadata.annotations}'
  # ИЛИ
  kubectl describe pod my-pod
  ```

---

### **7. Популярные аннотации в Kubernetes**

|Аннотация|Описание|
| --------------------| -----------------------------------------------------------------------------------------------------------|
|​`kubernetes.io/change-cause`​|Причина последнего обновления (используется в `kubectl rollout history`​)|
|​`kubectl.kubernetes.io/last-applied-configuration`​|JSON-конфигурация при последнем `kubectl apply`​ (для сравнения изменений)|
|​`sidecar.istio.io/inject: "true"`​|Включение автоматического внедрения Istio sidecar-контейнера|
|​`traefik.ingress.kubernetes.io/router-entrypoints: "websecure"`​|Настройка Traefik Ingress-контроллера|

---

### **8. Когда использовать аннотации?**

* **Для инструментов**: Если сторонний инструмент (например, Istio, Cert-Manager) требует настройки через аннотации.
* **Для документации**: Чтобы добавить пояснения к объектам (например, `description: "This service handles user authentication"`​).
* **Для CI/CD**: Хранение информации о сборке (`git-commit: "a1b2c3d"`​).

⚠ **Не используйте аннотации для данных, которые нужны для логики работы Kubernetes** — для этого есть `spec`​ и `labels`​.

---

### **9. Ограничения**

* Аннотации **не индексируются** (в отличие от меток), поэтому фильтрация по ним невозможна.
* Слишком большие аннотации могут снизить производительность API-сервера.
