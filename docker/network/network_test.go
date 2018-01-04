package dockernetwork

import (
	"testing"

	"github.com/cloudstax/firecamp/db"
	"github.com/cloudstax/firecamp/dns"
	"github.com/cloudstax/firecamp/server"
)

func TestAddDelIP(t *testing.T) {
	dbIns := db.NewMemDB()
	mockDNS := dns.NewMockDNS()
	serverIns := server.NewLoopServer()
	mockServerInfo := server.NewMockServerInfo()
	netIns := NewServiceNetwork(dbIns, mockDNS, serverIns, mockServerInfo)
	netIns.SetIfname("lo")

	ip := "10.0.0.1"
	err := netIns.AddIP(ip)
	if err != nil {
		t.Fatalf("add ip %s to lo error %s", ip, err)
	}

	err = netIns.DeleteIP(ip)
	if err != nil {
		t.Fatalf("del ip %s from lo error %s", ip, err)
	}
}
