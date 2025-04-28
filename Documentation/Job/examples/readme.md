```shell
kubectl apply -f simple-job.yaml
```



#### **Команды для управления**

| Действие                              | Команда |
| ----------------------------------------------- | ---------------- |
| Создать Job                            | `kubectl apply -f job.yaml`               |
| Просмотреть активные Job   | `kubectl get jobs`               |
| Проверить статус               | `kubectl describe job/my-job`               |
| Просмотреть логи Pod'а        | `kubectl logs job/my-job`               |
| Удалить Job                            | `kubectl delete job/my-job`               |
| Принудительно завершить | `kubectl patch job/my-job -p '{"spec":{"activeDeadlineSeconds":5}}'`               |           |espace=my-ns`               | Переключиться на Namespace                     |