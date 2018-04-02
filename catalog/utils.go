package catalog

import (
	"fmt"
	"strings"

	"github.com/cloudstax/firecamp/common"
	"github.com/cloudstax/firecamp/dns"
	"github.com/cloudstax/firecamp/manage"
	"github.com/cloudstax/firecamp/utils"
)

const hostSep = ","

// CreateSysConfigFile creates the content for the sys.conf file.
// example:
//   PLATFORM=ecs
//   SERVICE_MEMBER=mycas-0.cluster-firecamp.com
func CreateSysConfigFile(platform string, memberDNSName string) *manage.ReplicaConfigFile {
	content := fmt.Sprintf(sysContent, platform, memberDNSName)
	return &manage.ReplicaConfigFile{
		FileName: SYS_FILE_NAME,
		FileMode: common.DefaultConfigFileMode,
		Content:  content,
	}
}

// GenServiceMemberHosts creates the hostname list of all service members.
// Note: this currently works for the service that all members have a single name format, such as ZooKeeper, Kafka, etc.
// For service such as ElasticSearch, MongoDB, that has different name formats, this is not suitable.
func GenServiceMemberHosts(cluster string, service string, replicas int64) string {
	hosts := ""
	domain := dns.GenDefaultDomainName(cluster)
	for i := int64(0); i < replicas; i++ {
		member := utils.GenServiceMemberName(service, i)
		dnsname := dns.GenDNSName(member, domain)
		if len(hosts) == 0 {
			hosts = dnsname
		} else {
			hosts += hostSep + dnsname
		}
	}
	return hosts
}

// GenServiceMemberHostsWithPort creates the hostname:port list of all service members.
func GenServiceMemberHostsWithPort(attr *common.ServiceAttr, port int64) string {
	hosts := ""
	for i := int64(0); i < attr.Replicas; i++ {
		member := utils.GenServiceMemberName(attr.ServiceName, i)
		dnsname := dns.GenDNSName(member, attr.DomainName)
		if len(hosts) == 0 {
			hosts = fmt.Sprintf("%s:%d", dnsname, port)
		} else {
			hosts += fmt.Sprintf(",%s:%d", dnsname, port)
		}
	}
	return hosts
}

// TODO currently only support one jmx user.

// CreateJmxRemotePasswdConfFile creates the jmx remote password file.
func CreateJmxRemotePasswdConfFile(jmxUser string, jmxPasswd string) *manage.ReplicaConfigFile {
	return &manage.ReplicaConfigFile{
		FileName: JmxRemotePasswdConfFileName,
		FileMode: JmxConfFileMode,
		Content:  CreateJmxPasswdConfFileContent(jmxUser, jmxPasswd),
	}
}

// IsJmxPasswdConfFile checks if the file is jmx password conf file
func IsJmxPasswdConfFile(filename string) bool {
	return filename == JmxRemotePasswdConfFileName
}

// CreateJmxPasswdConfFileContent returns the jmxremote.password file content
func CreateJmxPasswdConfFileContent(jmxUser string, jmxPasswd string) string {
	return fmt.Sprintf(jmxPasswdFileContent, jmxUser, jmxPasswd)
}

// CreateJmxRemoteAccessConfFile creates the jmx remote access file.
func CreateJmxRemoteAccessConfFile(jmxUser string, accessPerm string) *manage.ReplicaConfigFile {
	return &manage.ReplicaConfigFile{
		FileName: JmxRemoteAccessConfFileName,
		FileMode: JmxConfFileMode,
		Content:  CreateJmxAccessConfFileContent(jmxUser, accessPerm),
	}
}

// IsJmxAccessConfFile checks if the file is jmx access conf file
func IsJmxAccessConfFile(filename string) bool {
	return filename == JmxRemoteAccessConfFileName
}

// CreateJmxAccessConfFileContent returns the jmxremote.access file content
func CreateJmxAccessConfFileContent(jmxUser string, accessPerm string) string {
	return fmt.Sprintf(jmxAccessFileContent, jmxUser, accessPerm)
}

// ReplaceJmxUserInAccessConfFile replaces the old jmx user in the access conf file
func ReplaceJmxUserInAccessConfFile(content string, newUser string, oldUser string) string {
	return strings.Replace(content, oldUser, newUser, 1)
}

// MBToBytes converts MB to bytes
func MBToBytes(mb int64) int64 {
	return mb * 1024 * 1024
}

const (
	separator  = "="
	sysContent = `
PLATFORM=%s
SERVICE_MEMBER=%s
`

	// jmx passwd file content format: "user passwd"
	jmxPasswdFileContent = "%s %s\n"

	// jmx access file content format: "user accessPermission"
	jmxAccessFileContent = "%s %s\n"
)
