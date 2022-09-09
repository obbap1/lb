package cgroups

import (
	"fmt"

	cgroups "github.com/containerd/cgroups"
)

func Start(filePath, mem string, port, minimum, maximum int, cpu float64) {
	// 2. Run Node Processes on random or specific ports

	// 3. Run Go server to forward requests via loopback interface or unix socket API

	// 4. Use control groups to limit resources and scale up if need be

	// 5. Scale down if need be as well

	// 6. Have metrics
	// res := cgroups.Resources{}
	// An empty string uses the default system slice
	// This is helpful not to create too many namespaces
	// the PID -1 creates a "parent cgroup"
	// The Cgroup is created with the file name
	// This means that the full path should be used when running a specific JS file
	// This might have nothing to do with JS and can be generic.
	// projectA/a.js and projectB/a.js will have two different namespaces.

	// Initially check if the CGroup exists

	m, err := cgroups.NewSystemd("user.slice", "my-cgroup-abc.slice", -1, &res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)

	control, err := cgroups.New(cgroups.Systemd, cgroups.Slice("user.slice"))

	// cgroups.
	// cgroup.MemoryThresholdEvent()
}
