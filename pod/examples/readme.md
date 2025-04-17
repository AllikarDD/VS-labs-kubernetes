```shell
kubectl apply -f simple-pod.yaml
```
```shell
kubectl apply -f simple-job.yaml
```

---

Чтобы получить список Pod'ов:

```bash
kubectl get pods
```

Подробная информация о Pod'е:

```bash
kubectl describe pod <имя-pod'а>
```

#### Удаление Pod'ов

Удаление Pod'а:

```bash
kubectl delete pod <имя-pod'а>
```

Однако если Pod управляется контроллером (например, Deployment), он будет автоматически воссоздан.

#### Доступ к Pod'ам

Для временного доступа к Pod'у можно использовать **port-forwarding**:

```bash
kubectl port-forward <имя-pod'а> 8080:80
```

Теперь можно обращаться к Pod'у через `localhost:8080`.

Для интерактивной отладки можно подключиться к контейнеру:

```bash
kubectl exec -it <имя-pod'а> -- /bin/bash
```

#### Логи Pod'ов

Просмотр логов:

```bash
kubectl logs <имя-pod'а>
```

Для Pod'ов с несколькими контейнерами укажите контейнер:

```bash
kubectl logs <имя-pod'а> -c <имя-контейнера>
```