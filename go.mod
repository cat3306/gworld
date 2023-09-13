module github.com/cat3306/goworld

go 1.18

require (
	github.com/alexeyco/simpletable v1.0.0
	github.com/cat3306/gocommon v0.0.0-20230911064234-07def0964423
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.3.0
	github.com/panjf2000/gnet/v2 v2.3.1
	github.com/sevlyar/go-daemon v0.1.6
	github.com/urfave/cli/v2 v2.25.7
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/zap v1.24.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.4
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/panjf2000/ants/v2 v2.4.8 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.etcd.io/etcd/api/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/v3 v3.5.9 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230303212802-e74f57abe488 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/panjf2000/gnet/v2 => github.com/cat3306/gnet/v2 v2.1.4-0.20230421080729-7e6031680b86
