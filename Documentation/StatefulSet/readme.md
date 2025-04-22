# StatefulSet (Набор с отслеживаниемым состоянием) в Kubernetes

Офицальная документация *StatefulSet*: [https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)

---

#### **1. Что такое StatefulSet?**

**StatefulSet** — это контроллер Kubernetes для управления **stateful-приложениями**, которые требуют:

* **Устойчивых идентификаторов** (постоянных имен Pod'ов и сетевых идентификаторов)
* **Упорядоченного развертывания и масштабирования** (последовательный запуск/остановка)
* **Устойчивого хранилища** (Persistent Volumes)

**Типичные use-cases**:

* Базы данных (MySQL, PostgreSQL, MongoDB)
* Кластерные приложения (ZooKeeper, Elasticsearch, etcd)
* Системы очередей (Kafka, RabbitMQ)

---

#### **2. Ключевые особенности**

|Характеристика|Описание|
| ------------------------------| ------------------------------------------------------------------------------------------------------------------------------------------|
|**Стабильные сетевые ID**|Каждый Pod получает постоянное имя (`<statefulset-name>-<ordinal>`​) и DNS-запись|
|**Упорядоченные операции**|Развертывание/удаление происходит последовательно (по номерам)|
|**Постоянное хранилище**|Каждому Pod'у выделяется свой PersistentVolume (при перезапуске данные сохраняются)|
|**Обновления**|Поддерживаются стратегии RollingUpdate и OnDelete|

---

#### **3. Пример StatefulSet для MySQL**

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: "mysql"  # Обязательно для работы сетевых идентификаторов
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:5.7
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
  volumeClaimTemplates:  # Динамическое создание PVC для каждого Pod'а
  - metadata:
      name: mysql-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 10Gi
```

---

#### **4. Жизненный цикл Pod'ов в StatefulSet**

1. **Создание**:

    * Pod'ы создаются последовательно (`mysql-0`​, `mysql-1`​, `mysql-2`​)
    * Каждый получает уникальный PersistentVolume
2. **Обновление**:

    * По умолчанию (`RollingUpdate`​) обновление идет в обратном порядке (`mysql-2`​ → `mysql-1`​ → `mysql-0`​)
3. **Удаление**:

    * Pod'ы удаляются в обратном порядке (`mysql-2`​ сначала, `mysql-0`​ последним)

---

#### **5. Сетевые идентификаторы**

Каждый Pod в StatefulSet получает:

* **Стабильное DNS-имя**:  
  ​`<pod-name>.<service-name>.<namespace>.svc.cluster.local`​  
  Пример: `mysql-0.mysql.default.svc.cluster.local`​
* **Headless Service** (обязателен для работы StatefulSet):

  ```yaml
  apiVersion: v1
  kind: Service
  metadata:
    name: mysql
  spec:
    clusterIP: None  # Headless Service
    ports:
    - port: 3306
    selector:
      app: mysql
  ```

---

#### **6. Команды для управления**

|Действие|Команда|
| ------------------------------------| ----------------|
|Создать StatefulSet|​`kubectl apply -f statefulset.yaml`​|
|Просмотреть StatefulSet|​`kubectl get statefulsets`​ (`sts`​)|
|Проверить Pod'ы|​`kubectl get pods -l app=mysql`​|
|Удалить StatefulSet|​`kubectl delete sts mysql`​|
|Масштабировать|​`kubectl scale sts mysql --replicas=5`​|
|Проверить тома|​`kubectl get pvc`​|

---

#### **7. Отличия от Deployment**

|Критерий|StatefulSet|Deployment|
| ------------------| -----------------------------------------------------| ----------------------------------------------------------|
|**Идентификаторы**|Стабильные имена (`app-0`​)|Динамические имена (`app-xyz`​)|
|**Хранилище**|Уникальное для каждого Pod'а|Общее (если не настроено иначе)|
|**Порядок операций**|Строгая последовательность|Параллельно|
|**Использование**|Stateful-приложения|Stateless-приложения|

---

#### **8. Ограничения**

* **Удаление StatefulSet не удаляет автоматически связанные PersistentVolumes** (нужно вручную очищать PVC)
* **Невозможно изменить volumeClaimTemplates после создания** (требуется пересоздание StatefulSet)
* **Более медленные операции** по сравнению с Deployment из-за последовательного выполнения

---

#### **9. Best Practices**

1. **Всегда используйте Headless Service** для сетевой идентификации.
2. **Настраивайте PodDisruptionBudget** для защиты от случайных остановок.
3. **Регулярно бэкапите PersistentVolumes** (особенно для баз данных).
4. **Для кластерных приложений** настраивайте корректные readinessProbes.

---
