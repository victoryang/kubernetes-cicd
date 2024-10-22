module github.com/victoryang/kubernetes-cicd

go 1.14

require (
	github.com/drone/drone-go v1.6.0
	github.com/drone/drone-yaml v1.2.3
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/google/go-github/v33 v33.0.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.1.1
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/xanzy/go-gitlab v0.40.2
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.9
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace k8s.io/client-go v9.0.0+incompatible => k8s.io/client-go v0.19.4
