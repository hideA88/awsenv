package version

import "github.com/hideA88/awsenv/pkg"

var (
	version   string
	gitCommit string
	build     string
)

func Show(ctx pkg.Context) {
	logger := ctx.Logger
	logger.Info("version=", version)
	logger.Info("gitCommit=", gitCommit)
	logger.Info("build=", build)
}
