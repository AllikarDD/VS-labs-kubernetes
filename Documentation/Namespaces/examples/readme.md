```shell
kubectl apply -f simple-namespace.yaml
```
```shell
kubectl apply -f simple-resourcequota.yaml
```


### **Команды для работы с Namespace**

| Команда | Описание                                              |
| ---------------- | --------------------------------------------------------------- |
| `kubectl get ns`               | Список всех Namespace                               |
| `kubectl create ns my-ns`               | Создать Namespace                                      |
| `kubectl delete ns my-ns`               | Удалить Namespace                                      |
| `kubectl get pods -n my-ns`               | Показать Pod'ы в определённом Namespace |
| `kubectl config set-context --current --namespace=my-ns`               | Переключиться на Namespace                     |