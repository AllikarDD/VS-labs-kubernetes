# CronJob (Периодические задания) в Kubernetes

Офицальная документация *CronJob*: [https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)

---

#### **1. Что такое CronJob?**

**CronJob** — это контроллер Kubernetes для запуска **заданий по расписанию**, аналогично классическому Unix-демону `cron`​. Основные сценарии использования:

* **Регулярные задачи обслуживания** (очистка логов, резервное копирование)
* **Периодические отчеты** (ежедневная аналитика, ночные ETL-процессы)
* **Планируемые операции** (массовые рассылки в определенное время)

**Как это работает**:  
CronJob создает объекты **Job** согласно заданному расписанию, которые затем выполняют свою работу через Pod'ы.

---

#### **2. Пример CronJob для резервного копирования**

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: nightly-backup
spec:
  schedule: "0 3 * * *"  # Каждый день в 3:00 (формат cron)
  jobTemplate:           # Шаблон для создания Job
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:14
            command: ["/bin/sh", "-c", "pg_dumpall -U postgres > /backups/$(date +%Y-%m-%d).sql"]
            volumeMounts:
            - name: backup-volume
              mountPath: /backups
          restartPolicy: OnFailure
          volumes:
          - name: backup-volume
            persistentVolumeClaim:
              claimName: backup-pvc
  successfulJobsHistoryLimit: 3  # Хранить историю 3 успешных Job
  failedJobsHistoryLimit: 1      # Хранить историю 1 неудачной Job
```

---

#### **3. Формат расписания (schedule)**

Использует стандартный **cron-формат** с 5 полями:

```
┌───────────── минуты (0-59)
│ ┌───────────── час (0-23)
│ │ ┌───────────── день месяца (1-31)
│ │ │ ┌───────────── месяц (1-12)
│ │ │ │ ┌───────────── день недели (0-6, 0=воскресенье)
│ │ │ │ │
│ │ │ │ │
* * * * *
```

**Примеры**:

* ​`"*/5 * * * *"`​ — каждые 5 минут
* ​`"0 22 * * 1-5"`​ — в 22:00 с понедельника по пятницу
* ​`"0 0 1 * *"`​ — в полночь первого числа каждого месяца

---

#### **4. Ключевые параметры**

|Параметр|Описание|
| ------------------| -----------------------------------------------------------------------------------------------------------------------------|
|**schedule**|Обязательное поле. Cron-строка расписания|
|**concurrencyPolicy**|​`Allow`​ (по умолчанию), `Forbid`​ (пропуск при наложении), `Replace`​ (замена текущей)|
|**startingDeadlineSeconds**|Максимальное время задержки запуска (если контроллер был выключен)|
|**successfulJobsHistoryLimit**|Сколько успешных Job хранить в истории (по умолчанию 3)|
|**failedJobsHistoryLimit**|Сколько неудачных Job хранить (по умолчанию 1)|
|**timeZone**|Часовой пояс для расписания (требуется Kubernetes 1.24+)|

---

#### **5. Команды для управления**

|Действие|Команда|
| -----------------------------------------------| ----------------|
|Создать CronJob|​`kubectl apply -f cronjob.yaml`​|
|Просмотреть CronJob|​`kubectl get cronjobs`​ (`cj`​)|
|Проверить статус|​`kubectl describe cronjob/nightly-backup`​|
|Просмотреть созданные Job|​`kubectl get jobs --watch`​|
|Вручную запустить Job|​`kubectl create job --from=cronjob/nightly-backup manual-run`​|
|Приостановить CronJob|​`kubectl patch cronjob/nightly-backup -p '{"spec":{"suspend":true}}'`​|
|Удалить CronJob|​`kubectl delete cronjob/nightly-backup`​|

---

#### **6. Политики выполнения (concurrencyPolicy)**

```yaml
spec:
  concurrencyPolicy: Replace  # Варианты: Allow, Forbid, Replace
```

* **Allow** (по умолчанию): Параллельное выполнение, если предыдущая Job еще работает
* **Forbid**: Пропуск нового запуска, если предыдущая Job активна
* **Replace**: Замена текущей Job на новую

---

#### **7. Best Practices**

1. **Используйте** **​`restartPolicy: OnFailure`​**​ для задач с возможностью рестарта
2. **Ограничивайте историю Job** для экономии ресурсов кластера
3. **Настраивайте** **​`startingDeadlineSeconds`​**​ для критичных по времени задач
4. **Логируйте результаты** в PersistentVolume или внешние системы
5. **Для коротких задач**

‍
