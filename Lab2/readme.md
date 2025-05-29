# Lab2

1. Для начала склонируйте репозитори https://github.com/AllikarDD/VS-labs-kubernetes  
   ‍

```bash
https://github.com/AllikarDD/VS-labs-kubernetes.git
```

# Volumes (Тома)

1. В папке `Documentation/Volumes/examples`​ создайте Volumes

```shell
apply -f simple-volume.yaml
```

2. Зайдите в поду и посмотрите создалась ли папка

# Ephemeral Volumes (Временные тома)

1. В папке `Documentation/Ephemeral Volumes/examples`​ создайте Ephemeral Volumes

```shell
kubectl apply -f simple-ephemeral-volume.yaml
```

2. Зайдите в поду и посмотрите создалась ли папка

# Secrets (Секреты)

1. В папке `Documentation/Secrets/examples`​ создайте секрет

```shell
kubectl apply -f simple-secret.yaml
```

2. Создайте поду с контейнером nginx
3. Добавьте секреты в этот контейнер через Volume
4. Зайдите в этот контейнер
5. Посмотрите что будет лежать в папке с секртами

# ConfigMap

1. В папке `Documentation/ConfigMap/examples`​ создайте ConfigMap

```shell
kubectl apply -f simple-configmap.yaml
```

2. Зайдите в поду
3. Проверьте значения переменных окружения, которые мы добавили

# Использование ConfigMap в коде

1. В файле Lab2/main.go измените код как показано ниже

```go
import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    greeting := os.Getenv("GREETING")
    text := os.Getenv("TEXT")
    fmt.Fprintln(w, greeting, ", this is Me, Mario")
    fmt.Fprintln(w, text)
}
```

2. Для примения изменений нужно пересобрать образ `myhello`​

3. Создайте конфигмапы

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmapref-example-cm
data:
   GREETING: Hola
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
   name: configmapkeyref-example-cm
data:
   text: "But also they call me Bombardino Crocodilo"
```

4. Добавьте ConfigMap в контейнер, который мы создавали в файле  Lab2/deployment.yaml

```yaml
#         Пример подключения всех параметров из конфигмапы configmapref-example-cm
          envFrom:
             - configMapRef:
                  name: configmapref-example-cm
                  
#         Пример подключения параметра text из конфигмапы configmapkeyref-example-cm
          env:
             - name: TEXT
               valueFrom:
                  configMapKeyRef:
                     name: configmapkeyref-example-cm
                     key: text
```

5. Запустите поду и посмотрите что получилось
6. Попробуйте отредактировать ConfigMap и обновите конфигурацию с помощью команды apply
7. Посмотрите поменялось ли что-то в вашем поднятом приложении
8. Попробуйте удалить поду
9. После поднятия поды, посмотрите применилась ли конфигурация

10. Так же можно создавать файлы конфигурации из конфигмапы

Создайте конфигмапу
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
   name: configfile-example-cm
data:
   file.conf: |
      TEXT=I stronger then 100 Monkeys
```
Создайте volumes в контейнере и свяжите с конфигмапой configfile-example-cm в файле  Lab2/deployment.yaml
```yaml
          volumeMounts:
            - name: config-volume
              mountPath: /etc/config
          
      volumes:
        - name: config-volume
          configMap:
            name: configfile-example-cm
```
Поскольку в нашем приложении сразу стартует веб сервис, мы не можем выполнять команды в контейнере 

Для того чтобы проверить, что файл появился, добавьте в Lab2/deployment.yaml
```yaml
      initContainers:
        - name: init-create-dir
          image: alpine
#         Выводим в консоль содержимое файла, который создали из конфигмапы
          command: [ 'sh', '-c', 'sleep 5 && cat /etc/config/file.conf' ]
          volumeMounts:
            - name: config-volume
              mountPath: /etc/config
```
Чтобы посмотреть логи контейнера, выполните команду 
```bash
kubectl logs <pod-name> -c <init-container-name>
```