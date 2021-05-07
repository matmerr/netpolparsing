package netpol

import (
	"io/ioutil"
	"log"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

func loadNetpol(filename string) (*netv1.NetworkPolicy, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDecoder()
	object := &netv1.NetworkPolicy{}
	err = runtime.DecodeInto(decoder, dat, object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func overlappingRange(n1 int, n2 *int32, m1 int, m2 *int32) bool {
	// first range is n1 to n2
	// second range is m1 to m2
	if n2 == nil || m2 == nil {
		return false
	}
	return int(*n2) <= m1 || n1 >= int(*m2)
}

func determineOverlap(ingressPorts, egressPorts map[netv1.NetworkPolicyPort]int) []netv1.NetworkPolicyPort {

	overlappedPorts := []netv1.NetworkPolicyPort{}

	for port := range ingressPorts {
		log.Printf("%+v:%+v", port.Port, *port.Protocol)
		for eg := range egressPorts {

			if port.Port.IntValue() == eg.Port.IntValue() && *port.Protocol == *eg.Protocol {
				log.Println(true)
			}

			if *port.Port == *eg.Port && *port.Protocol == *eg.Protocol { //&& *port.EndPort == *eg.EndPort {
				overlappedPorts = append(overlappedPorts, port)
				continue
			}

			// check if IP is in port range
			if overlappingRange(port.Port.IntValue(), port.EndPort, eg.Port.IntValue(), eg.EndPort) {
				overlappedPorts = append(overlappedPorts, port)
				continue
			}
		}
	}
	return overlappedPorts
}

func getPorts(netpol1 *netv1.NetworkPolicy) (map[netv1.NetworkPolicyPort]int, map[netv1.NetworkPolicyPort]int) {
	ingressPorts := make(map[netv1.NetworkPolicyPort]int)
	egressPorts := make(map[netv1.NetworkPolicyPort]int)

	for _, ingress := range netpol1.Spec.Ingress {
		for _, port := range ingress.Ports {
			ingressPorts[port] = 1
		}
	}

	for _, egress := range netpol1.Spec.Egress {
		for _, port := range egress.Ports {
			egressPorts[port] = 1
		}
	}

	return ingressPorts, egressPorts
}
