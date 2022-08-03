package nfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	defaultLogFile = "/tmp/ganesha.log"
	defaultPidFile = "/var/run/ganesha.pid"
)

var defaultConfig = []byte(`
NFS_Core_Param
{
    NLM_Port = 0;
    MNT_Port = 0;
    RQUOTA_Port = 0;
    Enable_NLM = false;
    Enable_RQUOTA = false;
    Enable_UDP = false;   
    fsid_device = false;
    Protocols = 4;
}

# uncomment to enable debug logging
# LOG { COMPONENTS { NFS_V4 = FULL_DEBUG; } }

NFSV4
{
    Lease_Lifetime = 60;
    Graceless = true;
    Grace_Period = 0;
    Minor_Versions = 1, 2;
    RecoveryBackend = fs_ng;
    RecoveryRoot = /tmp;
    Only_Numeric_Owners = true;
}

Export_defaults
{
    Protocols = 4;
    Transports = TCP;
    Access_Type = None;
    SecType = sys;
    Squash = None;
}

# Pseudo export, ganesha will automatically create one
# if one is not present
#EXPORT
#{
#    Export_Id = 0;
#    Protocols = 4;
#    Transports = TCP;
#    Access_Type = RW;
#    SecType = sys;
#    Squash = None;
#    Path = /export;
#    Pseudo = /;
#    FSAL { Name = VFS; }
#}
`)

type Server struct {
	logger     logrus.FieldLogger
	configPath string
	exportPath string
	exporter   *exporter
}

func NewServer(logger logrus.FieldLogger, configPath, exportPath, volume string) (*Server, error) {
	if err := setRlimitNOFILE(logger); err != nil {
		logger.Warnf("Error setting RLIMIT_NOFILE, there may be 'Too many open files' errors later: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err = ioutil.WriteFile(configPath, defaultConfig, 0600); err != nil {
			return nil, fmt.Errorf("error writing nfs config %s: %v", configPath, err)
		}
	}

	exporter, err := newExporter(logger, configPath, exportPath)
	if err != nil {
		return nil, fmt.Errorf("error creating nfs exporter: %v", err)
	}

	if _, err := exporter.CreateExport(volume); err != nil {
		return nil, err
	}

	return &Server{
		logger:     logger,
		configPath: configPath,
		exportPath: exportPath,
		exporter:   exporter,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	// Start ganesha.nfsd
	s.logger.Info("Running NFS server!")
	cmd := exec.CommandContext(ctx, "ganesha.nfsd", "-F", "-L", defaultLogFile, "-p", defaultPidFile, "-f", s.configPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ganesha.nfsd failed with error: %v, output: %s", err, out)
	}

	return nil
}

func setRlimitNOFILE(logger logrus.FieldLogger) error {
	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return fmt.Errorf("error getting RLIMIT_NOFILE: %v", err)
	}
	logger.Infof("starting RLIMIT_NOFILE rlimit.Cur %d, rlimit.Max %d", rlimit.Cur, rlimit.Max)
	rlimit.Max = 1024 * 1024
	rlimit.Cur = 1024 * 1024
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return err
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		return fmt.Errorf("error getting RLIMIT_NOFILE: %v", err)
	}
	logger.Infof("ending RLIMIT_NOFILE rlimit.Cur %d, rlimit.Max %d", rlimit.Cur, rlimit.Max)
	return nil
}
