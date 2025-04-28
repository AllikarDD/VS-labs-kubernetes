```shell
kubectl apply -f simple-deployment.yaml
```
```shell
kubectl apply -f simple-rollingupdate.yaml
```


### **Практические команды**

| Действие                                | Команда |
| ------------------------------------------------- | ---------------- |
| Проверить развертывание   | `kubectl get deployments`               |
| Показать детали                   | `kubectl describe deployment nginx-deployment`               |
| Обновить образ                     | `kubectl set image deployment/nginx-deployment nginx=nginx:1.19`               |
| Приостановить обновление | `kubectl rollout pause deployment/nginx-deployment`               |
| Возобновить обновление     | `kubectl rollout resume deployment/nginx-deployment`               |espace=my-ns`               | Переключиться на Namespace                     |