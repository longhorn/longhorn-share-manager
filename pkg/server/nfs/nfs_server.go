package nfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/guelfey/go.dbus"
	"github.com/sirupsen/logrus"
)

const (
	defaultExportBase  = "/tmp"
	defaultLogFile     = "/export/ganesha.log"
	defaultPidFile     = "/var/run/ganesha.pid"
	defaultConfigFile  = "/export/vfs.conf"
	defaultGracePeriod = 90
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

# enable debug logging for now, to test the node failover case
LOG { COMPONENTS { NFS_V4 = FULL_DEBUG; } }

NFSV4
{
    Lease_Lifetime = 10;
    Grace_Period = 90;
    Minor_Versions = 1, 2;
    RecoveryBackend = fs_ng;
}

Export_defaults
{
    Protocols = 4;
    Transports = TCP;
    Access_Type = None;
    SecType = sys;
    Squash = None;
}

# required pseudo export which creates the nfsv4 namespace
# ganesha creates a pseudo export internally if one is not specified
#EXPORT
#{
#    Export_Id = 0;
#    Filesystem_id = 0.0;
#    Protocols = 4;
#    Transports = TCP;
#    Access_Type = MDONLY_RO;
#    SecType = sys;
#    Squash = None;
#    Path = /export;
#    Pseudo = /;
#    FSAL { Name = PSEUDO; }
#}
`)

type Server struct {
	logger      logrus.FieldLogger
	configPath  string
	exportPath  string
	gracePeriod uint
	exporter    *exporter
}

func NewDefaultServer(logger logrus.FieldLogger) (*Server, error) {
	return NewServer(logger, defaultConfigFile, defaultGracePeriod)
}

func NewServer(logger logrus.FieldLogger, configPath string, gracePeriod uint) (*Server, error) {
	// Start dbus, needed for dynamic exports
	cmd := exec.Command("dbus-daemon", "--system")
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("dbus-daemon failed with error: %v, output: %s", err, out)
	}

	if err := setRlimitNOFILE(logger); err != nil {
		logger.Warnf("Error setting RLIMIT_NOFILE, there may be 'Too many open files' errors later: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err = ioutil.WriteFile(configPath, defaultConfig, 0600); err != nil {
			return nil, fmt.Errorf("error writing nfs config %s: %v", configPath, err)
		}
	}

	if err := setGracePeriod(configPath, gracePeriod); err != nil {
		return nil, fmt.Errorf("error writing grace period setting to nfs config: %v", err)
	}

	exporter, err := newExporter(logger, configPath, defaultExportBase)
	if err != nil {
		return nil, fmt.Errorf("error creating nfs exporter: %v", err)
	}

	return &Server{
		logger:      logger,
		configPath:  configPath,
		exportPath:  defaultExportBase,
		gracePeriod: gracePeriod,
		exporter:    exporter,
	}, nil
}

// Run : run the NFS NFSServer in the foreground until it exits
// Ideally, it should never exit when run in foreground mode
// We force foreground to allow the provisioner process to restart
// the NFSServer if it crashes - daemonization prevents us from using Wait()
// for this purpose
func (s *Server) Run(ctx context.Context) error {
	// Start ganesha.nfsd
	s.logger.Info("Running NFS server!")
	cmd := exec.CommandContext(ctx, "ganesha.nfsd", "-F", "-L", defaultLogFile, "-p", defaultPidFile, "-f", s.configPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ganesha.nfsd failed with error: %v, output: %s", err, out)
	}

	return nil
}

// Stop stops the nfs server.
func (s *Server) Stop() {
	// /bin/dbus-send --system   --dest=org.ganesha.nfsd --type=method_call /org/ganesha/nfsd/admin org.ganesha.nfsd.admin.shutdown
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

func setGracePeriod(ganeshaConfig string, gracePeriod uint) error {
	if gracePeriod > 180 {
		return fmt.Errorf("grace period cannot be greater than 180")
	}

	newLine := fmt.Sprintf("Grace_Period = %d;", gracePeriod)

	re := regexp.MustCompile("Grace_Period = [0-9]+;")

	read, err := ioutil.ReadFile(ganeshaConfig)
	if err != nil {
		return err
	}

	oldLine := re.Find(read)

	var file *os.File
	if oldLine == nil {
		// Grace_Period line not there, append the whole NFSV4 block.
		file, err = os.OpenFile(ganeshaConfig, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer file.Close()

		block := `NFSV4
		{
			Lease_Lifetime = 10;
			` + newLine + `
			Minor_Versions = 1, 2;
			RecoveryBackend = fs_ng;
		}`

		if _, err = file.WriteString(block); err != nil {
			return err
		}
		if err = file.Sync(); err != nil {
			return err
		}
	} else {
		// Grace_Period line there, just replace it
		replaced := strings.Replace(string(read), string(oldLine), newLine, -1)
		err = ioutil.WriteFile(ganeshaConfig, []byte(replaced), 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetExport(volume string) uint16 {
	return s.exporter.GetExport(volume)
}

func (s Server) GetExports() map[string]uint16 {
	return s.exporter.GetExportMap().volumeToid
}

func (s *Server) CreateExport(volume string) (uint16, error) {
	if id := s.exporter.GetExport(volume); id != 0 {
		return id, nil
	}

	id, err := s.exporter.CreateExport(volume)
	if err != nil {
		return 0, err
	}

	// TODO: we only need to run the dbus command if the server is running
	// export is persisted, let's update the running ganesha
	// ganesha dbus interface documentation
	// https://github.com/nfs-ganesha/nfs-ganesha/wiki/Dbusinterface
	err = func() error {
		volumePath := path.Join(s.exportPath, volume)
		conn, err := dbus.SystemBus()
		if err != nil {
			return fmt.Errorf("error getting dbus session bus: %v", err)
		}
		obj := conn.Object("org.ganesha.nfsd", "/org/ganesha/nfsd/ExportMgr")
		call := obj.Call("org.ganesha.nfsd.exportmgr.AddExport", 0, s.configPath, fmt.Sprintf("export(path = %s)", volumePath))
		if call.Err != nil {
			return fmt.Errorf("error calling org.ganesha.nfsd.exportmgr.AddExport: %v", call.Err)
		}
		return nil
	}()

	if err != nil {
		s.logger.WithField("volume", volume).WithError(err).Error("dbus export command failed")
		_ = s.exporter.DeleteExport(volume)
		return 0, fmt.Errorf("failed to export volume %v error: %v", volume, err)
	}

	return id, nil
}

func (s *Server) DeleteExport(volume string) (uint16, error) {
	id := s.exporter.GetExport(volume)
	if id == 0 {
		return id, nil
	}

	// TODO: it's possible we need to remove it first from the config
	// let's remove the export from the running ganesha then from the config
	// ganesha dbus interface documentation
	// https://github.com/nfs-ganesha/nfs-ganesha/wiki/Dbusinterface
	err := func() error {
		conn, err := dbus.SystemBus()
		if err != nil {
			return fmt.Errorf("error getting dbus session bus: %v", err)
		}
		obj := conn.Object("org.ganesha.nfsd", "/org/ganesha/nfsd/ExportMgr")
		call := obj.Call("org.ganesha.nfsd.exportmgr.RemoveExport", 0, id)
		if call.Err != nil {
			return fmt.Errorf("error calling org.ganesha.nfsd.exportmgr.RemoveExport: %v", call.Err)
		}
		return nil
	}()

	if err != nil {
		return id, fmt.Errorf("failed to delete export for volume %v error: %v", volume, err)
	}

	if err := s.exporter.DeleteExport(volume); err != nil {
		return id, err
	}

	return id, nil
}
