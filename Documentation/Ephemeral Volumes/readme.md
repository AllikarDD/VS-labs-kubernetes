# Ephemeral Volumes (Временные тома) в Kubernetes

Офицальная документация *Ephemeral Volumes*: [https://kubernetes.io/docs/concepts/storage/ephemeral-volumes/](https://kubernetes.io/docs/concepts/storage/ephemeral-volumes/)

---

#### **1. Что такое временные тома?**

**Временные тома** — это тома с жизненным циклом, привязанным к **отдельному Pod'у**. Они идеально подходят для:

* **Временных данных** (кэши, промежуточные файлы)
* **Сценариев, где данные не требуют сохранения** после завершения работы Pod'а
* **Работы с чувствительными данными**, которые должны быть удалены сразу после использования

**Ключевая особенность**:  
Данные удаляются при удалении Pod'а, но могут сохраняться между перезапусками контейнеров в рамках одного Pod'а.

---

#### **2. Основные типы временных томов**

|Тип тома|Описание|Преимущества|
| -----------------| ------------------------------------------------------------------------------------------------| -------------------------------------------------------------------------------|
|**emptyDir**|Пустой том, создается при запуске Pod'а|Простота настройки, общий для контейнеров|
|**generic ephemeral**|Упрощенная версия PersistentVolumeClaim для временных данных|Поддержка StorageClass, квот|
|**configMap/secret**|Специальные тома для конфигурации|Автоматическое обновление данных|
|**CSI ephemeral**|Временные тома от внешних CSI-драйверов|Поддержка сторонних решений|

---

#### **3. Примеры конфигураций**

#### **emptyDir (базовый временный том)**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-empty-dir
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: cache
      mountPath: /cache
  volumes:
  - name: cache
    emptyDir:
      sizeLimit: 500Mi  # Опциональное ограничение размера
```

#### **Generic Ephemeral Volume (динамическое временное хранилище)**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-gephem
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - name: data
      mountPath: "/data"
  volumes:
  - name: data
    ephemeral:
      volumeClaimTemplate:
        spec:
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: 1Gi
          storageClassName: "fast-storage"  # Использует StorageClass
```

---

#### **4. Сравнение с Persistent Volumes**

|Характеристика|Временные тома|Постоянные тома (PV/PVC)|
| ------------------------------| ---------------------------------------------------| ----------------------------------------------------------------------|
|**Жизненный цикл**|Привязан к Pod'у|Независим от Pod'ов|
|**Данные после удаления Pod'а**|Удаляются|Сохраняются|
|**Использование**|Кэши, временные вычисления|Базы данных, пользовательские данные|
|**Производительность**|Обычно выше|Зависит от бэкенда|

---

#### **5. Особые сценарии использования**

#### **Чувствительные данные**

```yaml
volumes:
- name: secret-data
  emptyDir:
    medium: Memory  # Хранится только в RAM (tmpfs)
```

#### **Общий кэш между контейнерами**

```yaml
containers:
- name: writer
  image: writer-app
  volumeMounts:
  - name: shared-cache
    mountPath: /cache

- name: reader
  image: reader-app
  volumeMounts:
  - name: shared-cache
    mountPath: /data

volumes:
- name: shared-cache
  emptyDir: {}
```

---

#### **6. Команды для управления**

|Действие|Команда|
| ------------------------------------------------------| ----------------|
|Проверить подключенные тома|​`kubectl describe pod/my-pod \| grep -A 10 Volumes`​|
|Создать Pod с временным томом|​`kubectl apply -f ephemeral-pod.yaml`​|
|Проверить использование|​`kubectl exec my-pod -- df -h`​|

---

#### **7. Best Practices**

1. **Для чувствительных данных** используйте `emptyDir.medium: Memory`​
2. **Устанавливайте sizeLimit** для предотвращения переполнения
3. **Для временных PV-подобных томов** используйте generic ephemeral volumes
4. **Не храните важные данные** — всегда считайте эти тома временными
5. **Для распределенных кэшей** рассмотрите отдельные решения (Redis, Memcached)

---

#### **8. Ограничения**

* **Нет гарантии сохранности данных** при пересоздании Pod'а
* **Нет встроенного шифрования** для emptyDir
* **Generic ephemeral volumes** требуют поддержки StorageClass
* **Максимальный размер** ограничен доступными ресурсами узла

---
