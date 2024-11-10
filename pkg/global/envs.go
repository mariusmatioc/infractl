package global

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func envToString(key string, val string) string {
	/* Need: "environment": [{"name": "env-name","value": "env-value"},...]	*/
	return fmt.Sprintf(`{"name": "%s", "value": "%s"}`, key, val)
}

// loadComposeEnvs loads envs from the docker compose file, with special processing if an RDS database is present (db)
func (serv *Service) loadComposeEnvs(composeFileFolder string) error {
	for key, val := range serv.composeConfig.Environment {
		if val == nil {
			serv.envMap[key] = ""
		} else {
			serv.envMap[key] = *val
		}
	}

	// The env files are nor processed. If needed, they should be in the craft file
	/*
		wd, _ := os.Getwd()
		defer func() { _ = os.Chdir(wd) }()
		err := os.Chdir(composeFileFolder) // we want relative file names to refer to compose folder
		if err != nil {
			return err
		}
		for _, name := range serv.composeConfig.EnvFile {
			err := serv.ReadEnvFile(name)
			if err != nil {
				return err
			}
		}
	*/
	return nil
}

//func (srv *Service) updateConnectedServiceEndpoint(serviceMap map[string]*Service) {
//	for key, val := range srv.envMap {
//		if _, ok := serviceMap[val]; ok {
//			// Another service. Map to load balancer
//			srv.envMap[key] = fmt.Sprintf("${aws_lb.%s.dns_name}", val)
//			break
//		}
//	}
//}

//func (srv *Service) updateDbHostEnvs(rds map[string]Rdb) {
//	for key, val := range srv.envMap {
//		if _, ok := rds[val]; ok {
//			// Maps to DB
//			srv.envMap[key] = fmt.Sprintf("${aws_db_instance.%s.address}", val)
//			break
//		}
//	}
//}
//
//func (srv *Service) updateMqHostEnvs(mqa map[string]Mq) {
//	for key, val := range srv.envMap {
//		if _, ok := mqa[val]; ok {
//			// Maps to DB
//
//			srv.envMap[key] = fmt.Sprintf("${aws_mq_broker.%s.instances.0.endpoints.0}", val)
//			break
//		}
//	}
//}

// updateCrossServiceHostEnvs updates the envs that refer to other services
// Another ECS, a DB or an MQ
func (srv *Service) updateCrossServiceHostEnvs(dependentServiceNames map[string]bool, endPointTemplate string) (err error) {
	// For envs in our service
	foundEnv := ""
	for key, val := range srv.envMap {
		if _, ok := dependentServiceNames[val]; ok {
			if foundEnv != "" {
				err = fmt.Errorf("service %s has more than one env that maps to service %s: %s %s", srv.Name, val, key, foundEnv)
				return
			}

			// This envs value maps to a dependent service name
			// Replace with endpoint
			srv.envMap[key] = fmt.Sprintf(endPointTemplate, val)
			foundEnv = key
		}
	}
	return
}

// GetEnvMapFromFiles assumes that relative names refer to root parent
func GetEnvMapFromFiles(files []string) (envMap map[string]string, err error) {
	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return
	}
	envMap = make(map[string]string)

	for _, name := range files {
		name = os.ExpandEnv(name)
		path, _ := filepath.Abs(name)
		err = ReadEnvFile(path, envMap)
		if err != nil {
			err = fmt.Errorf("error reading env file %s: %s", path, err.Error())
		}
	}
	return
}

// ReadEnvFile reads a file of format name=value and adds to envMap
func ReadEnvFile(name string, envMap map[string]string) error {
	text, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(text)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.TrimSpace(txt)
		if len(txt) == 0 || txt[0] == '#' {
			// Comment or empty line
			continue
		}
		// Extract name and value
		name := ""
		val := ""
		if i := strings.Index(txt, "="); i != -1 {
			name = txt[:i]
			val = txt[i+1:]
			var done bool
			val, done = RemoveQuotes(val)
			if !done {
				// Not quoted, discard comment
				if i = strings.Index(val, "#"); i != -1 {
					val = val[:i]
				}
			}
		} else {
			name = txt
		}
		name = strings.TrimSpace(name)
		val = strings.TrimSpace(val)
		envMap[name] = val
	}
	return nil
}

// ReadEnvFile reads a file of format name=value
func (srv *Service) ReadEnvFile(name string) error {
	return ReadEnvFile(name, srv.envMap)
}
