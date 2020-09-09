COMMIT_SHA1=$(git rev-parse --short HEAD || echo "0.0.0")
BUILD_TIME=$(date "+%F %T")
GO_VERSION=$(go version| awk '{print $3}')
go build -o stress -ldflags "-X github.com/oldthreefeng/stress/cmd.Version=$1 -X  github.com/oldthreefeng/stress/cmd.GoVersion=${GO_VERSION} -X github.com/oldthreefeng/stress/cmd.Build=${COMMIT_SHA1} -X 'github.com/oldthreefeng/stress/cmd.BuildTime=${BUILD_TIME}'" main.go && upx stress
