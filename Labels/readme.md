# Метки (Labels) и Селекторы (Selectors) в Kubernetes

**Истоничики**:

Офицальная документация *Labels and Selectors*: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/

---

### **1. Что такое метки (Labels)?**

**Метки** — это пары **ключ-значение**, которые присваиваются объектам Kubernetes (например, Pod'ам, Deployments, Services) для их организации и группировки.

**Примеры меток:**

* ​`environment: production`​
* ​`app: frontend`​
* ​`tier: backend`​
* ​`release: stable`​

Метки помогают:

* Управлять объектами Kubernetes.
* Фильтровать и находить нужные ресурсы.
* Определять, какие Pod'ы должны обрабатываться Service'ами, ReplicaSets и другими компонентами.

---

### **2. Синтаксис меток**

* **Ключи** могут состоять из:

  * Префикса (опционально) и имени.
  * Префикс (если есть) должен быть действительным DNS-поддоменом (например, `company.com/app`​).
  * Имя (обязательно) должно быть не длиннее 63 символов и может содержать `[a-z0-9A-Z]`​, `-`​, `_`​, `.`​.
* **Значения** должны быть не длиннее 63 символов (до 253 для префиксов) и могут содержать те же символы, что и ключи.

**Примеры допустимых меток:**

```yaml
metadata:
  labels:
    app: nginx
    tier: frontend
    environment: production
```

---

### **3. Селекторы (Selectors)**

Селекторы используются для выбора объектов по их меткам. Kubernetes поддерживает два типа селекторов:

#### **а) Селекторы на основе равенства (Equality-based)**

Фильтруют объекты по точному совпадению ключа и значения.  
Операторы: `=`​, `==`​, `!=`​.

**Примеры:**

* ​`environment = production`​ (выбирает объекты с `environment: production`​)
* ​`tier != frontend`​ (исключает объекты с `tier: frontend`​)

#### **б) Селекторы на основе множества (Set-based)**

Позволяют фильтровать объекты по набору значений.  
Операторы: `in`​, `notin`​, `exists`​.

**Примеры:**

* ​`environment in (production, staging)`​
* ​`tier notin (frontend, backend)`​
* ​`!app`​ (выбирает объекты без метки `app`​)

---

### **4. Использование меток и селекторов**

#### **Пример Deployment с метками и селекторами**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        tier: frontend
    spec:
      containers:
      - name: nginx
        image: nginx:latest
```

* ​**​`metadata.labels`​**​ — метки самого Deployment.
* ​**​`spec.selector.matchLabels`​**​ — определяет, какие Pod'ы управляются этим Deployment.
* ​**​`template.metadata.labels`​**​ — метки, которые будут присвоены создаваемым Pod'ам.

---

### **5. Применение меток**

Метки используются в:

* **Services** — для выбора Pod'ов, которые должны получать трафик.
* **ReplicaSets / Deployments** — для управления группами Pod'ов.
* **Node Selection** — для назначения Pod'ов на определённые узлы (`nodeSelector`​).
* **Администрировании** — для логической группировки (например, `team: devops`​, `project: alpha`​).

---

### **6. Команды kubectl для работы с метками**

* **Добавить метку к Pod:**

  ```sh
  kubectl label pods my-pod environment=production
  ```
* **Обновить метку:**

  ```sh
  kubectl label pods my-pod environment=staging --overwrite
  ```
* **Удалить метку:**

  ```sh
  kubectl label pods my-pod environment-
  ```
* **Фильтрация по меткам:**

  ```sh
  kubectl get pods -l environment=production
  kubectl get pods -l 'environment in (production, staging)'
  ```

---
