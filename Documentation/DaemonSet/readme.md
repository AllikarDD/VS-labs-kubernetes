# DaemonSet в Kubernetes

Офицальная документация *DaemonSet*:  [https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)

---

#### **1. Что такое DaemonSet?**

**DaemonSet** — это контроллер Kubernetes, который обеспечивает запуск **одного экземпляра Pod'а на каждом узле** (или на подмножестве узлов) кластера. Основные сценарии использования:

* **Системные демоны** (логирование, мониторинг, сбор метрик)
* **Сетевые плагины** (CNI, прокси, VPN)
* **Хранилище** (агенты для работы с дисками)

**Аналогия**:  
Как системные службы в Linux (например, `sshd`​ или `cron`​), но для Kubernetes-узлов.

---

#### **2. Ключевые особенности**

|Характеристика|Описание|
| ------------------------------| ---------------------------------------------------------------------------------------------------------------------------|
|**Один Pod на узел**|Автоматически развертывается на всех подходящих узлах|
|**Автоматическое добавление**|При добавлении нового узла Pod создается без ручного вмешательства|
|**Удаление вместе с узлом**|При удалении узла Pod также удаляется|
|**Использование ресурсов**|Обычно работает с `hostNetwork`​, `hostPID`​ или привилегированными контейнерами|

---

#### **3. Пример DaemonSet для fluentd (сбор логов)**

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd
  template:
    metadata:
      labels:
        name: fluentd
    spec:
      tolerations:  # Игнорирует ограничения узлов
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: fluentd
        image: fluent/fluentd:v1.14
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
```

---

#### **4. Размещение на подмножестве узлов**

DaemonSet можно ограничить для работы только на определенных узлах:

1. **Метки узлов** + **nodeSelector**:

    ```yaml
    spec:
      template:
        spec:
          nodeSelector:
            disktype: ssd  # Только узлы с этой меткой
    ```
2. **Tolerations** (для работы на master-узлах):

    ```yaml
    tolerations:
    - key: node-role.kubernetes.io/master
      effect: NoSchedule
    ```

---

#### **5. Команды для управления**

|Действие|Команда|
| ----------------------------------| ----------------|
|Создать DaemonSet|​`kubectl apply -f daemonset.yaml`​|
|Просмотреть DaemonSet|​`kubectl get daemonsets`​ (`ds`​)|
|Проверить Pod'ы|​`kubectl get pods -l name=fluentd -o wide`​|
|Удалить DaemonSet|​`kubectl delete ds fluentd`​|
|Проверить узлы|​`kubectl describe ds fluentd \| grep Nodes`​|

---

#### **6. Отличия от Deployment/ReplicaSet**

|Критерий|DaemonSet|Deployment/ReplicaSet|
| ------------------| -------------------------------------------------| -------------------------------------------------------|
|**Цель**|Системные задачи на узлах|Пользовательские приложения|
|**Масштабирование**|По количеству узлов|По произвольному replicas|
|**Обновления**|RollingUpdate (по умолчанию)|Полный набор стратегий|
|**Хранилище**|Часто использует hostPath|Обычно PersistentVolume|

---

#### **7. Обновление стратегии**

```yaml
spec:
  updateStrategy:
    type: RollingUpdate  # Варианты: RollingUpdate (по умолчанию) или OnDelete
    rollingUpdate:
      maxUnavailable: 1  # Максимум недоступных Pod'ов при обновлении
```

---

#### **8. Ограничения**

* **Нет горизонтального масштабирования** (количество Pod'ов равно количеству узлов)
* **Ограниченные стратегии обновления** по сравнению с Deployment
* **Часто требует привилегированного доступа** к узлу

---

#### **9. Best Practices**

1. **Используйте** **​`nodeSelector`​**​ **/**​**​`tolerations`​**​ для контроля размещения.
2. **Ограничивайте ресурсы** (CPU/memory) для системных демонов.
3. **Проверяйте readinessProbe** для критически важных демонов.
4. **Для обновлений** предпочитайте `RollingUpdate`​ с `maxUnavailable=1`​.

---
