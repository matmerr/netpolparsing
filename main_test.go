package netpol

import (
	"log"
	"testing"
)

func TestDecode(t *testing.T) {

	netpol1, err := loadNetpol("./egressblockport25.yaml")
	if err != nil {
		t.Fatal(err)
	}

	//netpol2, err := loadNetpol("smtp_relay_policy.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ingressPorts, egressPorts := getPorts(netpol1)


	overlap := determineOverlap(ingressPorts, egressPorts)


	log.Printf("overlap: %+v", overlap)
	t.FailNow()
}
