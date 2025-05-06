# Volumes (Тома) в Kubernetes

Офицальная документация *Volumes*: [https://kubernetes.io/docs/concepts/storage/volumes/](https://kubernetes.io/docs/concepts/storage/volumes/)

---

#### **1. Что такое Volume в Kubernetes?**

**Volume** — это механизм для предоставления **постоянного хранилища** контейнерам в Pod'ах. В отличие от локального хранилища контейнера, томы:

* **Сохраняют данные** после перезапуска контейнера
* **Могут быть общими** для нескольких контейнеров в Pod'е
* **Поддерживают различные** бэкенды (локальные диски, облачные хранилища, NFS и др.)

**Ключевая особенность**:  
Тома привязаны к **жизненному циклу Pod'а**, но не к контейнеру.

---

#### **2. Основные типы томов**

|Тип тома|Описание|Пример использования|
| -----------------| -------------------------------------------------------------------------------------------------| ----------------------------------------------------------------------|
|**emptyDir**|Временное хранилище, создается при запуске Pod'а|Кэши, временные файлы|
|**hostPath**|Доступ к файлам на узле (ноде)|Доступ к docker.sock, системным логам|
|**PersistentVolume (PV)**|Внешнее хранилище, управляемое администратором|Базы данных, пользовательские данные|
|**ConfigMap/Secret**|Специальные тома для конфигураций и секретов|Настройки приложений, TLS-сертификаты|
|**Cloud Provider Volumes**|Интеграция с облачными хранилищами (AWS EBS, GCE PD, Azure Disk)|Облачные приложения|

---

#### **3. Примеры конфигураций**

#### **emptyDir (временное хранилище)**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pd
spec:
  containers:
  - image: nginx
    name: nginx-container
    volumeMounts:
    - mountPath: /cache
      name: cache-volume
  volumes:
  - name: cache-volume
    emptyDir: {}
```

#### **hostPath (доступ к файлам узла)**

```yaml
volumes:
- name: hostpath-volume
  hostPath:
    path: /data
    type: Directory  # Варианты: File, Directory, Socket и др.
```

#### **PersistentVolumeClaim** 

```yaml
volumes:
- name: pvc-volume
  persistentVolumeClaim:
    claimName: my-pvc  # Должен существовать PVC
```

---

#### **4. Подключение томов к контейнерам**

Тома подключаются через **volumeMounts** в спецификации контейнера:

```yaml
containers:
- name: app
  image: nginx
  volumeMounts:
  - name: data-volume  # Должен соответствовать имени в volumes
    mountPath: "/usr/share/nginx/html"
    readOnly: true     # Опционально
```

---

#### **5. Работа с Persistent Volumes (PV) и PVC**

1. **PersistentVolume (PV)**  — ресурс хранилища в кластере (создается администратором)
2. **PersistentVolumeClaim (PVC)**  — запрос на выделение хранилища (создается пользователем)

**Пример PVC**:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard  # Опционально
```

---

#### **6. Команды для управления**

|Действие|Команда|
| ----------------------------------------| ----------------|
|Просмотреть тома Pod'а|​`kubectl describe pod/my-pod \| grep -A 10 Volumes`​|
|Создать PVC|​`kubectl apply -f pvc.yaml`​|
|Просмотреть PV/PVC|​`kubectl get pv`​, `kubectl get pvc`​|
|Удалить PVC|​`kubectl delete pvc/my-pvc`​|

---

#### **7. Best Practices**

1. **Для production-сред** всегда используйте **PersistentVolume** вместо hostPath
2. **Ограничивайте доступ** для sensitive томов (`readOnly: true`​)
3. **Проверяйте reclaimPolicy** для PV (`Retain`​/`Delete`​/`Recycle`​)
4. **Используйте StorageClass** для динамического выделения хранилища
5. **Для конфигов** предпочитайте **ConfigMap/Secret** над ручным монтированием файлов

---

#### **8. Ограничения**

* **emptyDir** удаляется при удалении Pod'а
* **hostPath** создает риски безопасности (доступ к узлу)
* **PV/PVC** требуют предварительной настройки администратором
* **Миграция данных** между узлами не всегда автоматическая
