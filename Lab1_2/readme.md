# Lab1_2

1. Для начала склонируйте репозитори [https://github.com/AllikarDD/VS-lab1-kubernetes.git](https://github.com/AllikarDD/VS-lab1-kubernetes.git) :

‍

```bash
https://github.com/AllikarDD/VS-lab1-kubernetes.git
```

# Пространства имён (Namespaces)

1. В папке `Documentation/Namespaces/examples`​ создайте Namespace

```shell
kubectl apply -f simple-namespace.yaml
```

2. Создайте ограничения на ресурсы в вашем Namespace

```yaml
kubectl apply -f simple-resourcequota.yaml
```

3. Проверьте информацию о вашем Namespace

4. Удалите Namespace

Команды для Namespace расположены в файле `Documentation/Namespaces/examples/readme.md`​

# Job (Задание)

1. В папке `Documentation/Job/examples`​ создайте Job

```shell
kubectl apply -f simple-job.yaml
```

2. Job создаёт поду, которая считает число Пи до 2000 знаков после запятой, и выводит это число в логи
3. Выведите логи этой поды
4. Попробуйте  запустить job еще раз, что будет с прошлой подой?

Команды для Job расположены в файле `Documentation/Job/examples/readme.md`​

# DaemonSet

1. В папке `Documentation/DaemonSet/examples`​ создайте DaemonSet

```shell
kubectl apply -f simple-daemonset.yaml
```

2. Проверьте подробную информацию о DaemonSet

Команды для DaemonSet расположены в файле `Documentation/DaemonSet/examples/readme.md`​

# CronJob (Периодические задания)

1. В папке `Documentation/CronJob/examples`​ создайте CronJob

```shell
kubectl apply -f simple-cronjob.yaml
```

2. Посмотрите список CronJob
3. Попробуйте запустить CronJob вручную

Команды для CronJob расположены в файле `Documentation/CronJob /examples/readme.md`​

# StatefulSet (Набор с отслеживаниемым состоянием)

1. В папке `Documentation/StatefulSet/examples`​ создайте StatefulSet

```shell
kubectl apply -f simple-statefulset.yaml
```

2. Посмотрите запустилась ли пода
3. Попробуйте решить проблему и запустить её

Команды для StatefulSet расположены в файле `Documentation/StatefulSet/examples/readme.md`​

# ReplicaSet (Набор реплик)

1. В папке `Documentation/ReplicaSet/examples`​ создайте ReplicaSet

```shell
kubectl apply -f simple-replicaset.yaml
```

2. Посмотрите информацию о ReplicaSet
3. Увеличьте кол-во реплик
4. Посмотрите сколько стало под
5. Попробуйте решить проблему и запустить её

Команды для ReplicaSet расположены в файле `Documentation/ReplicaSet/examples/readme.md`​

# Deployment (Развертывание)

1. Соберем образ myhello

```bash
docker image build -t myhello .
```

2. Для проверки запустим контейнер

```bash
docker container run -p 9999:8888 myhello
```

В результате у вас должна заработать ваша собственная копия демонстрационного приложения. Чтобы в этом убедиться, откройте  URL-адрес (http://localhost:9999/).

3. Запуск демонстрационного приложения

```bash
kubectl run demo --image=myhello--port=9999 --labels app=demo --image-pull-policy='Never'
```

*Напомним, что вам следует перенаправить порт 9999 вашего локального компьютера к порту 8888 контейнера, чтобы к нему можно было подключаться с помощью веб-браузера. То же самое необходимо сделать здесь, используя команду kubectl port-forward:*

```bash
kubectl port-forward demo 9999:8888
```

В результате у вас должна заработать ваша собственная копия демонстрационного приложения. Чтобы в этом убедиться, откройте  URL-адрес (http://localhost:9999/)

## Запуск приложения с помощью манифеста

В качестве примера возьмем наше демонстрационное приложение (deployment.yaml):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo
  labels:
    app: demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      labels:
        app: demo
    spec:
      containers:
        - name: demo
          image: myhello
          imagePullPolicy: Never
          ports:
          - containerPort: 8888
```

В папке *​`Lab1_2`​*​используйте команду:

```yaml
kubectl apply -f deployment.yaml
```

Проверим что пода создалась

```yaml
kubectl get pods
или 
kubectl get pods --selector app=demo
```

Должны увидеть вот такую картину

![](assets/image-20250217021810-p258qef.png)

Попробуем удалить поду

```bash
kubectl delete pods --selector app=demo
```

Если еще раз посмотрите на поду то она удалилась, появилась новая с новым названием

![](assets/image-20250217022200-414e7oz.png)

Это произошло, потому что сработало автоматическое развертывание. Чтобы удалить поду, чтобы она не восстанавливалась, нужно удалить deployment командой

```bash
kubectl delete deployment demo
```

Попробуйте поменять, и применить изменения

```yaml
spec:
  replicas: 1 -> 2
```

Вы увидите что стало две реплики вашего приложения

![](assets/image-20250217022313-iwkv858.png)

‍
