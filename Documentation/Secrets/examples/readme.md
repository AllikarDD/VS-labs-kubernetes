```shell
kubectl apply -f simple-secret.yaml
```


#### **Команды для управления**

##### Просмотр списка
`kubectl get secrets`

##### Детали секрета (без значений)
`kubectl describe secret db-creds`

#### Получение значений (декодирование)
`kubectl get secret db-creds -o jsonpath='{.data.password}' | base64 --decode     |`