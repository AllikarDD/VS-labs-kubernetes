```shell
kubectl apply -f simple-replicaset.yaml
```


### **Основные команды**

| Действие                | Команда |
| --------------------------------- | ---------------- |
| Создать ReplicaSet       | `kubectl apply -f replicaset.yaml`               |
| Просмотреть          | `kubectl get rs`               |
| Проверить детали | `kubectl describe rs/frontend`               |
| Удалить                  | `kubectl delete rs/frontend`               |
| Масштабирова