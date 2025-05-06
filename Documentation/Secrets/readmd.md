# Secrets (Секреты) в Kubernetes

Офицальная документация *Secrets*: [https://kubernetes.io/docs/concepts/configuration/secret/](https://kubernetes.io/docs/concepts/configuration/secret/)

---

#### **1. Что такое Secret?**

**Secret** — это объект Kubernetes для безопасного хранения конфиденциальных данных:

* Пароли и токены
* TLS-сертификаты
* SSH-ключи
* Ключи API

**Ключевые особенности**:

* Хранение в закодированном виде (base64)
* Ограниченный доступ через RBAC
* Отделение секретов от образов контейнеров

---

#### **2. Типы Secrets**

|Тип|Использование|Пример|
| --------| ------------------------------------------------------------------------| ---------------------------------|
|**Opaque**|Произвольные пользовательские данные|Пароли БД|
|**kubernetes.io/tls**|TLS-сертификаты|HTTPS-терминация|
|**docker-registry**|Учетные данные для Docker Registry|Доступ к registry|
|**service-account-token**|Токены ServiceAccount|Авторизация Pod'ов|

---

#### **3. Создание Secret**

#### **Через kubectl**

```bash
# Из литералов
kubectl create secret generic db-creds \
  --from-literal=username=admin \
  --from-literal=password='S!B\*d$zDsb='

# Из файлов
kubectl create secret generic tls-cert \
  --from-file=tls.crt=./cert.pem \
  --from-file=tls.key=./key.pem
```

#### **Через YAML (закодированный в base64)**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=  # echo -n "admin" | base64
  password: UyFCXCpkJHpEc2I9
```

---

#### **4. Использование Secrets в Pod'ах**

#### **Как переменные среды**

```yaml
containers:
- name: app
  image: my-app
  env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-creds
          key: password
```

#### **Как файлы через Volume**

```yaml
volumes:
- name: secret-volume
  secret:
    secretName: tls-cert

volumeMounts:
- name: secret-volume
  mountPath: "/etc/ssl"
  readOnly: true
```

---

#### **5. Безопасность Secrets**

#### **Рекомендации**:

1. **Шифрование на rest**: Включите EncryptionConfig
2. **RBAC**: Ограничьте доступ по принципу наименьших привилегий
3. **Вращение секретов**: Регулярно обновляйте
4. **Интеграция с внешними хранилищами**:

    * HashiCorp Vault
    * AWS Secrets Manager
    * Azure Key Vault

#### **Ограничения**:

* Base64 ≠ Шифрование
* Секреты доступны всем с доступом к etcd
* По умолчанию сохраняются в памяти узла

---

#### **6. Пример: TLS Secret для Ingress**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ingress-tls
type: kubernetes.io/tls
data:
  tls.crt: LS0t...tCg==
  tls.key: LS0t...LQo=
```

Использование в Ingress:

```yaml
spec:
  tls:
  - hosts:
    - example.com
    secretName: ingress-tls
```

---

#### **7. Best Practices**

1. **Не коммитьте** секреты в Git
2. **Используйте сторонние инструменты** для управления:

    * Sealed Secrets
    * External Secrets Operator
3. **Регулярно ротируйте** (особенно после увольнений)
4. **Аудит доступа** через Kubernetes audit logging
5. **Ограничьте mount** только необходимым Pod'ам

---

#### **8. Мониторинг и управление**

```bash
# Просмотр списка
kubectl get secrets

# Детали секрета (без значений)
kubectl describe secret db-creds

# Получение значений (декодирование)
kubectl get secret db-creds -o jsonpath='{.data.password}' | base64 --decode
```

---
