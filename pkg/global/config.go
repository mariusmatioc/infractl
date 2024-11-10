package global

// getManagedServices returns all database as well as queue managed services and removes them from the services list
// servs are the remaining generic services
func (config *Config) getManagedServices(services Services) (servs Services, rdsMap map[string]Rdb, mqsMap map[string]Mq, err error) {
	rdsMap, servs, err = config.GetRds(services)
	if err != nil {
		return
	}
	mqsMap, servs, err = config.GetMqs(servs)
	return
}

func (cfg *Config) GetPortsFor(name string) (ports map[int]bool) {
	ports = make(map[int]bool)
	for _, svc := range cfg.Compose.Services {
		if svc.Name == name {
			for _, port := range svc.Ports {
				ports[int(port.Target)] = true
			}
			break
		}
	}
	return
}
