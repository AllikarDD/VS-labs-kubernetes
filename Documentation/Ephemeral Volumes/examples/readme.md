```shell
kubectl apply -f simple-ephemeral-volume.yaml
```


#### **Команды для управления**

| Действие                                     | Команда |
| ------------------------------------------------------ | ---------------- |
| Проверить подключенные тома | `kubectl describe pod/my-pod \| grep -A 10 Volumes`               |
| Создать Pod с временным томом  | `kubectl apply -f ephemeral-pod.yaml`               |
| Проверить использование        | `kubectl exec my-pod -- df -h`               |c/my-pvc`               |e64 --decode     |`