package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"math/rand"
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
	Config    Config
	Client    *api.Client
	ServiceID string
}

func New(params Config) (*Discovery, error) {
	cfg := api.DefaultConfig()
	cfg.Address = params.Address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Discovery{
		Config:    params,
		Client:    client,
		ServiceID: fmt.Sprintf("%s-%d", params.ServiceCluster, params.ServiceId),
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

func (d *Discovery) Deregister() error {
	return d.Client.Agent().ServiceDeregister(d.ServiceID)
}

func (d *Discovery) GetService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := d.Client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	return services, nil

}

func (d *Discovery) GetServiceAddress(serviceName string) (string, error) {
	services, err := d.GetService(serviceName)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("service %s not found", serviceName)
	}

	idx := rand.Intn(len(services))
	service := services[idx]

	address := service.Service.Address
	if address == "" {
		address = service.Node.Address
	}

	return fmt.Sprintf("%s:%d", address, service.Service.Port), nil
}
