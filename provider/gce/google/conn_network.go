// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package google

import (
	"github.com/juju/errors"

	"github.com/juju/juju/network"
)

// Ports build a list of all open port ranges for a given firewall name
// (within the Connection's project) and returns it. If the firewall
// does not exist then the list will be empty and no error is returned.
func (gce Connection) Ports(fwname string) ([]network.PortRange, error) {
	firewall, err := gce.raw.GetFirewall(gce.ProjectID, fwname)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Annotate(err, "while getting ports from GCE")
	}

	var ports []network.PortRange
	for _, allowed := range firewall.Allowed {
		for _, portRangeStr := range allowed.Ports {
			portRange, err := network.ParsePortRange(portRangeStr)
			if err != nil {
				return ports, errors.Annotate(err, "bad ports from GCE")
			}
			portRange.Protocol = allowed.IPProtocol
			ports = append(ports, *portRange)
		}
	}

	return ports, nil
}

// OpenPorts sends a request to the GCE API to open the provided port
// ranges on the named firewall. If the firewall does not exist yet it
// is created, with the provided port ranges opened. Otherwise the
// existing firewall is updated to add the provided port ranges to the
// ports it already has open. The call blocks until the ports are
// opened or the request fails.
func (gce Connection) OpenPorts(fwname string, ports []network.PortRange) error {
	// TODO(ericsnow) Short-circuit if ports is empty.

	// Compose the full set of open ports.
	currentPorts, err := gce.Ports(fwname)
	if err != nil {
		return errors.Trace(err)
	}
	inputPortsSet := network.NewPortSet(ports...)
	if inputPortsSet.IsEmpty() {
		return nil
	}
	currentPortsSet := network.NewPortSet(currentPorts...)

	// Send the request, depending on the current ports.
	if currentPortsSet.IsEmpty() {
		firewall := firewallSpec(fwname, inputPortsSet)
		if err := gce.raw.AddFirewall(gce.ProjectID, firewall); err != nil {
			return errors.Annotatef(err, "opening port(s) %+v", ports)
		}

	} else {
		newPortsSet := currentPortsSet.Union(inputPortsSet)
		firewall := firewallSpec(fwname, newPortsSet)
		if err := gce.raw.UpdateFirewall(gce.ProjectID, fwname, firewall); err != nil {
			return errors.Annotatef(err, "opening port(s) %+v", ports)
		}
	}
	return nil
}

// ClosePorts sends a request to the GCE API to close the provided port
// ranges on the named firewall. If the firewall does not exist nothing
// happens. If the firewall is left with no ports then it is removed.
// Otherwise it will be left with just the open ports it has that do not
// match the provided port ranges. The call blocks until the ports are
// closed or the request fails.
func (gce Connection) ClosePorts(fwname string, ports []network.PortRange) error {
	// Compose the full set of open ports.
	currentPorts, err := gce.Ports(fwname)
	if err != nil {
		return errors.Trace(err)
	}
	inputPortsSet := network.NewPortSet(ports...)
	if inputPortsSet.IsEmpty() {
		return nil
	}
	currentPortsSet := network.NewPortSet(currentPorts...)
	newPortsSet := currentPortsSet.Difference(inputPortsSet)

	// Send the request, depending on the current ports.
	if newPortsSet.IsEmpty() {
		// TODO(ericsnow) Handle case where firewall does not exist.
		if err := gce.raw.RemoveFirewall(gce.ProjectID, fwname); err != nil {
			return errors.Annotatef(err, "closing port(s) %+v", ports)
		}
	} else {
		firewall := firewallSpec(fwname, newPortsSet)
		if err := gce.raw.UpdateFirewall(gce.ProjectID, fwname, firewall); err != nil {
			return errors.Annotatef(err, "closing port(s) %+v", ports)
		}
	}
	return nil
}
