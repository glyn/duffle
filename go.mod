module github.com/cnabio/duffle

go 1.13

require (
	github.com/Masterminds/semver v1.5.0
	github.com/cnabio/cnab-go v0.8.2-beta1
	github.com/containerd/cgroups v0.0.0-20200108155730-918ed86e29cc // indirect
	github.com/docker/cli v0.0.0-20191017083524-a8ff7f821017
	github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible
	github.com/docker/docker v1.4.2-0.20181229214054-f76d6a078d88
	github.com/docker/go v1.5.1-1
	github.com/docker/go-connections v0.4.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gosuri/uitable v0.0.4
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/oklog/ulid v1.3.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/pivotal/image-relocation v0.0.0-20191111101224-e94aff6df06c
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/technosophos/moniker v0.0.0-20180509230615-a5dbd03a2245
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
