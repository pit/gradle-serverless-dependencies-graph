# terraform-serverless-private-registry
AWS Serverless Terraform Private Registry


# S3 Bucket files structure

```
- /modules
-- /namespace
--- /name
---- /provider
----- /version

Sample:

- /modules


source = terraform.kvinta.io/apps/epc-normalizer/k8s
version = v1.12.1

-- /apps
--- /epc-normalizer
---- /k8s
----- /v1.0.0
----- /v1.0.1
----- /v1.2.1
----- /v1.3.1
----- /v1.12.1


source = terraform.kvinta.io/kvinta/application/k8s
version = v1.12.1

-- /kvinta
--- /application
---- /k8s
----- /v1.0.0
----- /v1.0.1
----- /v1.2.1
----- /v1.3.1
----- /v1.12.1


source = terraform.kvinta.io/kvinta/kubernetes/yandex
version = v1.12.1

-- /kvinta
--- /kubernetes
---- /yandex
----- /v1.0.0
----- /v1.0.1
----- /v1.2.1
----- /v1.3.1
----- /v1.12.1

source = terraform.kvinta.io/kvinta/eks/aws
version = v1.12.1

-- /kvinta
--- /eks
---- /aws
----- /v1.0.0
----- /v1.0.1
----- /v1.2.1
----- /v1.3.1
----- /v1.12.1

source = terraform.kvinta.io/terraform/kvinta-eks/aws
version = v1.12.1

-- /terraform
--- /k8s
---- /yandex
----- /v1.0.0
----- /v1.0.1
----- /v1.2.1
----- /v1.3.1
----- /v1.12.1


```