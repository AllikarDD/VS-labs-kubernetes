# Лабораторная работа 3

## Предоставление доступа к Pod в кластере

Создайте файл `run-my-nginx.yaml`​

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
spec:
  selector:
    matchLabels:
      run: my-nginx
  replicas: 2
  template:
    metadata:
      labels:
        run: my-nginx
    spec:
      containers:
      - name: my-nginx
        image: nginx
        ports:
        - containerPort: 80
```

Это делает Pod доступным с любого узла в вашем кластере. Проверьте узлы, на которых работают Pod:

```shell
kubectl apply -f ./run-my-nginx.yaml
kubectl get pods -l run=my-nginx -o wide
```

```
NAME                        READY     STATUS    RESTARTS   AGE       IP            NODE
my-nginx-3800858182-jr4a2   1/1       Running   0          13s       10.244.3.4    kubernetes-minion-905m
my-nginx-3800858182-kna2y   1/1       Running   0          13s       10.244.2.5    kubernetes-minion-ljyd
```

Проверьте IP-адреса ваших Pod:

```shell
kubectl get pods -l run=my-nginx -o custom-columns=POD_IP:.status.podIPs
    POD_IP
    [map[ip:10.244.3.4]]
    [map[ip:10.244.2.5]]
```

Вы должны иметь возможность подключиться по SSH к любому узлу в кластере и использовать инструменты типа `curl`​ для отправки запросов на оба IP-адреса. Обратите внимание, что контейнеры *не* используют порт 80 на узле, и нет никаких специальных NAT-правил для маршрутизации трафика к Pod. Это означает, что вы можете запускать несколько nginx Pod на одном узле, все использующие один и тот же `containerPort`​, и получать к ним доступ с любого другого Pod или узла в кластере, используя назначенный Pod IP-адрес. Если вы хотите настроить переадресацию с определенного порта на узле Node к Pod, вы можете это сделать - но модель сети Kubernetes предполагает, что в этом нет необходимости.

## Создание Service

Итак, у нас есть Pod, работающие с nginx в плоском, кластерном адресном пространстве. Теоретически, вы можете обращаться к этим Pod напрямую, но что произойдет, если узел выйдет из строя? Pod завершат работу вместе с ним, а ReplicaSet внутри Deployment создаст новые Pod с другими IP-адресами. Именно эту проблему решает Service.

Service в Kubernetes - это абстракция, которая определяет логическую группу Pod в вашем кластере, предоставляющих одинаковый функционал. При создании каждому Service назначается уникальный IP-адрес (также называемый clusterIP). Этот адрес связан с жизненным циклом Service и не изменится, пока Service существует. Pod могут быть настроены для взаимодействия с Service, зная, что связь с Service будет автоматически балансироваться между Pod, входящими в Service.

Вы можете создать Service для ваших 2 реплик nginx с помощью kubectl expose:

```yaml
kubectl expose deployment/my-nginx
```

Это эквивалентно применению следующего yaml с помощью kubectl apply -f:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-nginx
  labels:
    run: my-nginx
spec:
  ports:
  - port: 80
    protocol: TCP
  selector:
    run: my-nginx
```

Эта конфигурация создаст Service, который направляет трафик на TCP-порт 80 любого Pod с меткой `run: my-nginx`​, и предоставляет его на абстрактном порту Service (`targetPort`​ - это порт, который контейнер принимает трафик, `port`​ - это абстрактный порт Service, который может быть любым портом, используемым другими Pod для доступа к Service). Ознакомьтесь с API-объектом [Service](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#service-v1-core), чтобы увидеть список поддерживаемых полей в определении Service. Проверьте ваш Service:

```shell
kubectl get svc my-nginx
```

```
NAME       TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
my-nginx   ClusterIP   10.0.162.149   <none>        80/TCP    21s
```

Как упоминалось ранее, Service поддерживается группой Pod. Эти Pod доступны через [EndpointSlices](https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/). Селектор Service будет непрерывно оцениваться, а результаты будут отправляться в EndpointSlice, который связан с Service с помощью [меток](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels). Когда Pod завершает работу, он автоматически удаляется из EndpointSlices, которые содержат его как конечную точку. Новые Pod, соответствующие селектору Service, автоматически добавляются в EndpointSlice для этого Service. Проверьте конечные точки и обратите внимание, что IP-адреса совпадают с Pod, созданными на первом этапе:

```shell
kubectl describe svc my-nginx
```

```
Name:                my-nginx
Namespace:           default
Labels:              run=my-nginx
Annotations:         <none>
Selector:            run=my-nginx
Type:                ClusterIP
IP Family Policy:    SingleStack
IP Families:         IPv4
IP:                  10.0.162.149
IPs:                 10.0.162.149
Port:                <unset> 80/TCP
TargetPort:          80/TCP
Endpoints:           10.244.2.5:80,10.244.3.4:80
Session Affinity:    None
Events:              <none>
```

```shell
kubectl get endpointslices -l kubernetes.io/service-name=my-nginx
```

```
NAME             ADDRESSTYPE   PORTS   ENDPOINTS               AGE
my-nginx-7vzhx   IPv4          80      10.244.2.5,10.244.3.4   21s
```

Теперь вы должны иметь возможность отправлять curl-запросы к nginx Service по адресу `<CLUSTER-IP>:<PORT>`​ с любого узла в вашем кластере. Обратите внимание, что IP-адрес Service полностью виртуальный и никогда не передается по сети. Если вам интересно, как это работает, вы можете прочитать больше о [сервис-прокси](https://kubernetes.io/docs/reference/networking/virtual-ips/).

Для того чтобы проверить можете восползоваться командами:

```yaml
kubectl run tester --image=alpine -it --rm -- /bin/sh
# Внутри пода:
apk add curl
curl http://<CLUSTER-IP>
exit
```

### DNS

Kubernetes предоставляет дополнительный сервис DNS, который автоматически назначает DNS-имена другим сервисам. Вы можете проверить его работу в вашем кластере:

```shell
kubectl get services kube-dns --namespace=kube-system
```

```
NAME       TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)         AGE
kube-dns   ClusterIP   10.0.0.10    <none>        53/UDP,53/TCP   8m
```

В дальнейшем мы будем предполагать, что у вас есть сервис с постоянным IP-адресом (my-nginx) и DNS-сервер, который назначил имя этому IP. Здесь мы используем аддон CoreDNS (название приложения `kube-dns`​), поэтому вы можете обращаться к сервису из любого pod в кластере стандартными методами (например, `gethostbyname()`​). Если CoreDNS не работает, вы можете включить его, следуя инструкциям в [CoreDNS README](https://github.com/coredns/deployment/tree/master/kubernetes) или [Installing CoreDNS](https://kubernetes.io/docs/tasks/administer-cluster/coredns/#installing-coredns). Давайте запустим еще одно приложение curl для тестирования:

```yaml
kubectl run tester --image=alpine -it --rm -- /bin/sh
# Внутри пода:
apk add curl
nslookup my-nginx
exit
```

Затем нажмите Enter и выполните команду `nslookup my-nginx`​:

```shell
[ root@curl-131556218-9fnch:/ ]$ nslookup my-nginx
Сервер:    10.0.0.10
Адрес 1: 10.0.0.10

Имя:      my-nginx
Адрес 1: 10.0.162.149
```

## Защита сервиса

До этого момента мы обращались к nginx-серверу только изнутри кластера. Прежде чем открывать доступ к сервису из интернета, необходимо обеспечить безопасность канала связи. Для этого потребуется:

* Самоподписанные сертификаты для HTTPS (если у вас уже нет сертификата)
* Nginx-сервер, настроенный на использование сертификатов
* [Secret](https://kubernetes.io/docs/concepts/configuration/secret/), который делает сертификаты доступными для pod

Следуйте этим инструкциям:

```shell
# Создаем пару открытый/закрытый ключ
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout nginx.key -out nginx.crt -subj "/CN=my-nginx/O=my-nginx"
# Конвертируем ключи в base64
cat nginx.crt | base64 -w 0
cat nginx.key | base64 -w 0
```

Используйте вывод этих команд для создания yaml-файла. Закодированные в base64 значения должны быть в одной строке.

```yaml
apiVersion: "v1"
kind: "Secret"
metadata:
  name: "nginxsecret"
  namespace: "default"
type: kubernetes.io/tls
data:
# Вставить значения из cat nginx.crt | base64 -w 0
  tls.crt: ""
# Вставить значения из cat nginx.key | base64 -w 0
  tls.key: ""
```

Теперь создаем secret с помощью файла:

```shell
kubectl apply -f nginxsecrets.yaml
kubectl get secrets
```

```
NAME                  TYPE                                  DATA      AGE
nginxsecret           kubernetes.io/tls                     2         1m
```

А также configmap:

```shell
kubectl create configmap nginxconfigmap --from-file=default.conf
```

Пример файла `default.conf`​ можно найти в [репозитории примеров Kubernetes](https://github.com/kubernetes/examples/tree/bc9ca4ca32bb28762ef216386934bef20f1f9930/staging/https-nginx/).

```
configmap/nginxconfigmap created
```

```shell
kubectl get configmaps
```

```
NAME             DATA   AGE
nginxconfigmap   1      114s
```

Подробности ConfigMap `nginxconfigmap`​ можно посмотреть следующей командой:

```shell
kubectl describe configmap nginxconfigmap
```

Вывод будет примерно таким:

```console
Name:         nginxconfigmap
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
default.conf:
----
server {
        listen 80 default_server;
        listen [::]:80 default_server ipv6only=on;

        listen 443 ssl;

        root /usr/share/nginx/html;
        index index.html;

        server_name localhost;
        ssl_certificate /etc/nginx/ssl/tls.crt;
        ssl_certificate_key /etc/nginx/ssl/tls.key;

        location / {
                try_files $uri $uri/ =404;
        }
}

BinaryData
====

Events:  <none>
```

Теперь модифицируем реплики nginx для запуска HTTPS-сервера с использованием сертификатов из secret, а также Service для открытия обоих портов (80 и 443):

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-nginx
  labels:
    run: my-nginx
spec:
  type: NodePort
  ports:
    - port: 8080
      targetPort: 80
      protocol: TCP
      name: http
    - port: 443
      protocol: TCP
      name: https
  selector:
    run: my-nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-nginx
spec:
  selector:
    matchLabels:
      run: my-nginx
  replicas: 1
  template:
    metadata:
      labels:
        run: my-nginx
    spec:
      volumes:
        - name: secret-volume
          secret:
            secretName: nginxsecret
        - name: configmap-volume
          configMap:
            name: nginxconfigmap
      containers:
        - name: nginxhttps
          image: nginx
          ports:
            - containerPort: 443
            - containerPort: 80
          volumeMounts:
            - mountPath: /etc/nginx/ssl
              name: secret-volume
            - mountPath: /etc/nginx/conf.d
              name: configmap-volume
```

Важные моменты манифеста nginx-secure-app:

* Он содержит спецификации и Deployment, и Service в одном файле
* [Nginx-сервер](https://github.com/kubernetes/examples/tree/master/staging/https-nginx/default.conf) обслуживает HTTP-трафик на порту 80 и HTTPS-трафик на порту 443, а Service открывает оба порта
* Каждый контейнер получает доступ к ключам через том, смонтированный в `/etc/nginx/ssl`​. Это настраивается *до* запуска nginx-сервера

```shell
kubectl delete deployments,svc my-nginx; kubectl create -f ./nginx-secure-app.yaml
```

Теперь можно обратиться к nginx-серверу с любого узла.

```shell
kubectl get pods -l run=my-nginx -o custom-columns=POD_IP:.status.podIPs
    POD_IP
    [map[ip:10.244.3.5]]
```

```yaml
kubectl run tester --image=alpine -it --rm -- /bin/sh
# Внутри пода:
apk add curl
#Обратите внимание что вызываем https, а не http
curl -k https://<POD-IP>
exit
```

```shell
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLS connect error: error:00000000:lib(0)::reason(0)
* OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to 10.1.0.84:443
* closing connection #0
curl: (35) TLS connect error: error:00000000:lib(0)::reason(0)

```

Обратите внимание, что мы использовали параметр `-k`​ в curl на последнем шаге - это потому, что во время генерации сертификатов мы ничего не знали о pod, на которых будет работать nginx, поэтому приходится указывать curl игнорировать несоответствие CName. Создав Service, мы связали CName, использованный в сертификате, с фактическим DNS-именем, используемым pod при поиске Service. Протестируем это из pod (для простоты повторно используем тот же secret, pod нужен только nginx.crt для доступа к Service):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: curl-deployment
spec:
  selector:
    matchLabels:
      app: curlpod
  replicas: 1
  template:
    metadata:
      labels:
        app: curlpod
    spec:
      volumes:
        - name: secret-volume
          secret:
            secretName: nginxsecret
      containers:
        - name: curlpod
          command:
            - sh
            - -c
            - while true; do sleep 1; done
          image: alpine
          volumeMounts:
            - mountPath: /etc/nginx/ssl
              name: secret-volume
```

```shell
kubectl apply -f ./curlpod.yaml
kubectl get pods -l app=curlpod
```

```
NAME                               READY     STATUS    RESTARTS   AGE
curl-deployment-1515033274-1410r   1/1       Running   0          1m
```

```shell
kubectl exec -it curl-deployment-1515033274-1410r7 --  curl -v https://my-nginx  --cacert /etc/nginx/ssl/tls.crt 

```

## Публикация сервиса во внешнюю сеть

Для некоторых компонентов приложения может потребоваться предоставить доступ к сервису через внешний IP-адрес. Kubernetes поддерживает два основных способа публикации сервисов: NodePort и LoadBalancer. Созданный ранее сервис уже использует NodePort, поэтому ваши реплики nginx с HTTPS готовы принимать интернет-трафик, если узлы имеют публичные IP-адреса.

### Проверка настроек NodePort

```shell
kubectl get svc my-nginx -o yaml | grep nodePort -C 5
```

Вывод показывает:

```yaml
  uid: 07191fb3-f61a-11e5-8ae5-42010af00002
spec:
  clusterIP: 10.0.162.149
  ports:
  - name: http
    nodePort: 31704    # Порт для HTTP-трафика
    port: 8080
    protocol: TCP
    targetPort: 80
  - name: https
    nodePort: 32453    # Порт для HTTPS-трафика
    port: 443
    protocol: TCP
    targetPort: 443
  selector:
    run: my-nginx
```

### Поиск внешних IP-адресов узлов

```shell
kubectl get nodes -o yaml | grep ExternalIP -C 1
```

Пример вывода:

```yaml
    - address: 104.197.41.11
      type: ExternalIP
    allocatable:
--
    - address: 23.251.152.56
      type: ExternalIP
    allocatable:
```

### Тестирование доступа

```shell
curl https://<ВНЕШНИЙ-IP>:<ПОРТ-NODE> -k
```

В ответе должна быть страница nginx:

```html
<h1>Welcome to nginx!</h1>
```

## Настройка балансировщика нагрузки

Изменим тип сервиса с NodePort на LoadBalancer:

```shell
kubectl edit svc my-nginx
kubectl get svc my-nginx
```

Результат:

```
NAME       TYPE           CLUSTER-IP     EXTERNAL-IP        PORT(S)               AGE
my-nginx   LoadBalancer   10.0.162.149   xx.xxx.xxx.xxx     8080:30163/TCP        21s
```

Тестирование:

```shell
curl https://<ВНЕШНИЙ-IP> -k
```

Ответ:

```html
<title>Welcome to nginx!</title>
```

### Особенности работы:

1. **EXTERNAL-IP** - публичный адрес, доступный из интернета
2. **CLUSTER-IP** - внутренний адрес, работающий только внутри кластера
