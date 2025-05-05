package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type Config struct {
	Address        string
	ServiceIp      string
	ServicePort    int
	HealthURI      string
	ServiceCluster string
	ServiceId      int
}

type Discovery struct {
	Config Config
	Client *api.Client
}

func New(params Config) (*Discovery, error) {
	cfg := api.DefaultConfig()
	cfg.Address = params.Address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Discovery{
		Config: params,
		Client: client,
	}, nil
}

func (d *Discovery) Register() error {

	consulReg := &api.AgentServiceRegistration{
		ID:   fmt.Sprintf("%s-%d", d.Config.ServiceCluster, d.Config.ServiceId),
		Name: d.Config.ServiceCluster,
		Port: d.Config.ServicePort,
		Check: &api.AgentServiceCheck{
			Interval:                       (time.Second * 10).String(),
			Timeout:                        (time.Second * 5).String(),
			HTTP:                           fmt.Sprintf("http://%s:%d/%s", d.Config.ServiceIp, d.Config.ServicePort, d.Config.HealthURI),
			DeregisterCriticalServiceAfter: (time.Second * 10).String(),
		},
	}

	return d.Client.Agent().ServiceRegister(consulReg)
}
