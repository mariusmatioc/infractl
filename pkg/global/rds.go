package global

import (
	"fmt"
	"os"
	"strconv"
)

func (config *Config) GetRds(services Services) (rdsMap map[string]Rdb, servs Services, err error) {
	rec := config.GetEcsRecipe()
	servs = services

	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return
	}
	rdsMap = make(map[string]Rdb)
	for name, db := range rec.SimpleRds.Databases {
		rds := Rdb{
			Name:        name,
			DbEngine:    db.DbEngine,
			Public:      db.Public,
			MachineType: db.MachineType,
			StorageType: db.StorageType,
			StorageIops: db.StorageIops,
			StorageGigs: db.StorageGigs,
		}
		envMap := map[string]string{}
		for _, filePath := range db.EnvFiles {
			filePath = os.ExpandEnv(filePath)
			err = ReadEnvFile(filePath, envMap)
			if err != nil {
				return
			}
		}
		rds.UserName = envMap["DB_USER"]
		rds.Password = envMap["DB_PWD"]
		rds.DbName = envMap["DB_NAME"]
		rds.Port, err = strconv.Atoi(envMap["DB_PORT"])
		if err != nil {
			err = fmt.Errorf("%s: DB_PORT is not an integer", name)
			return
		}
		msg := rds.Validate()
		if msg != "" {
			err = fmt.Errorf("invalid parameter for '%s' in simpl_rds: %s", name, msg)
			return
		}
		rdsMap[name] = rds
	}

	// Now remove any RDs from the general services list
	if len(rdsMap) == 0 {
		return
	}
	// Remove services that are to be converted to RDS
	removed := 0
	for i := len(servs) - 1; i >= 0; i-- {
		svcName := servs[i].Name
		if rds, ok := rdsMap[svcName]; ok {
			if rds.UserName == "" {
				userName := servs[i].envMap["DB_USER"]
				if userName == "" {
					err = fmt.Errorf("service %s has no DB_USER defined", svcName)
					return
				}
				rds.UserName = userName
			}
			if rds.Password == "" {
				pwd := servs[i].envMap["DB_PWD"]
				if pwd == "" {
					err = fmt.Errorf("service %s has no DB_PWD defined", svcName)
					return
				}
				rds.Password = pwd
			}
			if rds.DbName == "" {
				dbName := servs[i].envMap["DB_NAME"]
				if dbName == "" {
					err = fmt.Errorf("service %s has no DB_NAME defined", svcName)
					return
				}
				rds.DbName = dbName
			}
			if rds.Port == 0 {
				port := servs[i].envMap["DB_PORT"]
				if port == "" {
					err = fmt.Errorf("service %s has no DB_PORT defined", svcName)
					return
				}
				rds.Port, err = strconv.Atoi(port)
				if err != nil {
					return
				}
			}
			// Remove
			servs = append(servs[:i], servs[i+1:]...)
			delete(config.NameToService, svcName)
			removed++
		}
	}
	if removed != len(rdsMap) {
		err = fmt.Errorf("not all simpl_rds services were found in docker compose")
		return
	}
	return
}

func (rd *Rdb) Validate() string {
	if rd.DbEngine == "" {
		return "engine is required"
	}
	if rd.DbName == "" {
		return "db_name is required"
	}
	if rd.Port == 0 {
		return "port is required"
	}
	if rd.MachineType == "" {
		return "machine_type is required"
	}
	if rd.StorageGigs == 0 {
		return "storage_gigs is required"
	}
	return ""
}
