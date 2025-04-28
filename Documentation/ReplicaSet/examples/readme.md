```shell
kubectl apply -f simple-statefulset.yaml
```


#### **Команды для управления**

| Действие                   | Команда |
| ------------------------------------ | ---------------- |
| Создать StatefulSet         | `kubectl apply -f statefulset.yaml`               |
| Просмотреть StatefulSet | `kubectl get statefulsets` (`sts`)            |
| Проверить Pod'ы          | `kubectl get pods -l app=mysql`               |
| Удалить StatefulSet         | `kubectl delete sts mysql`               |
| Масштабировать       | `kubectl scale sts mysql --replicas=5`               |
| Проверить тома        | `kubectl get pvc`               |l patch job/my-job -p '{"spec":{"activeDeadlineSeconds":5}}'`               |           |espace=my-ns`               | Переключиться на Namespace                     |