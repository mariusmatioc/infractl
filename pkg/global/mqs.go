package global

import (
	"fmt"
	"os"
)

// GetMqs returns the message queues defined in the simple_mqs section of the recipe and removes them from the services
func (config *Config) GetMqs(services Services) (mqsMap map[string]Mq, servs Services, err error) {
	rec := config.GetEcsRecipe()
	servs = services

	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	err = os.Chdir(RootParent) // we want relative file names to refer to parent of root folder
	if err != nil {
		return
	}
	mqsMap = make(map[string]Mq)
	for name, q := range rec.SimpleMqs.Queues {
		mq := Mq{
			Name:   name,
			Engine: q.Engine,
			Public: q.Public,
		}
		envs, err2 := GetEnvMapFromFiles(q.EnvFiles)
		if err2 != nil {
			err = err2
			return
		}
		needed := []string{"RABBITMQ_USER", "RABBITMQ_PWD"}
		for _, n := range needed {
			if envs[n] == "" {
				err = fmt.Errorf("service '%s' has no %s defined", name, needed)
				return
			}
		}
		mq.UserName = envs["RABBITMQ_USER"]
		mq.Password = envs["RABBITMQ_PWD"]
		msg := mq.validate()
		if msg != "" {
			err = fmt.Errorf("invalid parameter for '%s' in simple_mqs: %s", name, msg)
			return
		}
		mqsMap[name] = mq
	}

	// Now remove any messages queues from the general services list
	if len(mqsMap) == 0 {
		return
	}
	// Remove services that are to be converted to message queues
	removed := 0
	for i := len(servs) - 1; i >= 0; i-- {
		svcName := servs[i].Name
		if mq, ok := mqsMap[svcName]; ok {
			if mq.UserName == "" {
				userName := servs[i].envMap["RABBITMQ_DEFAULT_USER"]
				if userName == "" {
					err = fmt.Errorf("service %s has no RABBITMQ_DEFAULT_USER defined", svcName)
					return
				}
				mq.UserName = userName
			}
			if mq.Password == "" {
				pwd := servs[i].envMap["RABBITMQ_DEFAULT_PASS"]
				if pwd == "" {
					err = fmt.Errorf("service %s has no RABBITMQ_DEFAULT_PASS defined", svcName)
					return
				}
				mq.Password = pwd
			}
			// Remove
			servs = append(servs[:i], servs[i+1:]...)
			delete(config.NameToService, svcName)
			removed++
		}
	}
	if removed != len(mqsMap) {
		err = fmt.Errorf("not all simple_mqs services were found in docker compose")
		return
	}
	return
}

func (mq *Mq) validate() string {
	if mq.Engine != "RabbitMQ" {
		return "engine must be RabbitMQ"
	}
	return ""
}
