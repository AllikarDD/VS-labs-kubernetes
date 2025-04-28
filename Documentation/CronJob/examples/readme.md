```shell
kubectl apply -f simple-cronjob.yaml
```


#### **Команды для управления**

| Действие                              | Команда |
| ----------------------------------------------- | ---------------- |
| Создать CronJob                        | `kubectl apply -f cronjob.yaml`               |
| Просмотреть CronJob                | `kubectl get cronjobs` (`cj`)            |
| Проверить статус               | `kubectl describe cronjob/nightly-backup`               |
| Просмотреть созданные Job | `kubectl get jobs --watch`               |
| Вручную запустить Job         | `kubectl create job --from=cronjob/nightly-backup manual-run`               |
| Приостановить CronJob            | `kubectl patch cronjob/nightly-backup -p '{"spec":{"suspend":true}}'`               |
| Удалить CronJob                        | `kubectl delete cronjob/nightly-backup`               |