```shell
kubectl apply -f simple-volume.yaml
```


#### **Команды для управления**

| Действие                       | Команда |
| ---------------------------------------- | ---------------- |
| Просмотреть тома Pod'а | `kubectl describe pod/my-pod \| grep -A 10 Volumes`               |
| Создать PVC                     | `kubectl apply -f pvc.yaml`               |
| Просмотреть PV/PVC          | `kubectl get pv`, `kubectl get pvc`             |
| Удалить PVC                     | `kubectl delete pvc/my-pvc`               |e64 --decode     |`