```shell
kubectl apply -f simple-daemonset.yaml
```


#### **Команды для управления**

| Действие                 | Команда |
| ---------------------------------- | ---------------- |
| Создать DaemonSet         | `kubectl apply -f daemonset.yaml`               |
| Просмотреть DaemonSet | `kubectl get daemonsets` (`ds`)            |
| Проверить Pod'ы        | `kubectl get pods -l name=fluentd -o wide`               |
| Удалить DaemonSet         | `kubectl delete ds fluentd`               |
| Проверить узлы      | `kubectl describe ds fluentd \| grep Nodes`               |