package mikrotik

import (
	"fmt"
	"log"

	"gopkg.in/routeros.v2"
)

type mikrotikRepository struct {
	client *routeros.Client
}

type Mikrotik interface {
	GetIdentity() (map[string]string, error)
	GetSecrets(btsName string, btsHost string) ([]map[string]string, error)
	GetActiveConnections() ([]map[string]string, error)
	EnableSNMP()
	SetMacFromAC()
	SetRemoteAddress()
	GetAddressList(listName string) []map[string]string
	AddSecretToAddressList(ip string, comment string, addressList string)
	RemoveSecretFromAddressList(addressListData []map[string]string, addressListComment string)
}

func NewMikrotikRepository(addr, user, password string) (Mikrotik, error) {
	dial, err := routeros.Dial(addr+":8728", user, password)
	if err != nil {
		return nil, err
	}

	return &mikrotikRepository{client: dial}, nil
}

func (mr *mikrotikRepository) GetIdentity() (map[string]string, error) {
	identity := []map[string]string{}
	mkt, err := mr.client.Run("/system/identity/print")
	if err != nil {
		return nil, err
	}

	for _, r := range mkt.Re {
		identity = append(identity, r.Map)
	}
	return identity[0], nil
}

func (mr *mikrotikRepository) GetSecrets(
	btsName string,
	btsHost string,
) ([]map[string]string, error) {
	secret := []map[string]string{}

	mkt, err := mr.client.Run("/ppp/secret/print")
	if err != nil {
		return nil, err
	}

	for _, d := range mkt.Re {
		d.Map["host"] = btsHost
		d.Map["bts"] = btsName
		secret = append(secret, d.Map)
	}
	return secret, nil
}

func (mr *mikrotikRepository) GetActiveConnections() ([]map[string]string, error) {
	var activeUsers []map[string]string
	reply, err := mr.client.Run("/ppp/active/print")
	if err != nil {
		return nil, err
	}

	for _, r := range reply.Re {
		activeUsers = append(activeUsers, r.Map)
	}

	return activeUsers, nil
}

func (mr *mikrotikRepository) EnableSNMP() {
	reply, err := mr.client.Run("/snmp/set", "=enabled=yes", "=trap-version=2")
	if err != nil {
		log.Println(err)
	}
	log.Println(reply.Done.Word)
}

func (mr *mikrotikRepository) SetMacFromAC() {
	reply, err := mr.client.Run("/ppp/secret/print")
	if err != nil {
		log.Println("error to get secrets users", err)
	}
	var secret []map[string]string

	for _, s := range reply.Re {
		fmt.Println(s.Map)
		secret = append(secret, s.Map)
	}

	activeUsers, err := mr.GetActiveConnections()
	if err != nil {
		log.Println("error to get active connections", err)
	}

	for _, su := range secret {
		for _, au := range activeUsers {
			if su["name"] == au["name"] {
				_, err := mr.client.Run(
					"/ppp/secret/set",
					fmt.Sprintf("=numbers=%s", su["name"]),
					fmt.Sprintf("=caller-id=%s", au["caller-id"]),
				)
				if err != nil {
					log.Println("error to set mac", err)
				}

			}
		}
	}
}

func (mr *mikrotikRepository) SetRemoteAddress() {
	secrets, err := mr.GetSecrets("", "")
	if err != nil {
		log.Println(err)
	}
	activeConnections, err := mr.GetActiveConnections()
	if err != nil {
		log.Println("error to get active connections")
	}

	for _, secret := range secrets {
		for _, active := range activeConnections {
			if secret["name"] == active["name"] {
				_, err := mr.client.Run(
					"/ppp/secret/set",
					fmt.Sprintf("=numbers=%s", secret["name"]),
					fmt.Sprintf("=remote-address=%s", active["address"]),
				)
				if err != nil {
					log.Println("error to set remote address")
				}
			}
		}
	}
}

func (mr *mikrotikRepository) AddSecretToAddressList(ip string, comment string, addressList string) {

	reply, _ := mr.client.Run(
		"/ip/firewall/address-list/add",
		fmt.Sprintf("=list=%s", addressList),
		fmt.Sprintf("=address=%s", ip),
		fmt.Sprintf("=comment=%s", comment),
	)
	log.Println(reply)
}

func (mr *mikrotikRepository) RemoveSecretFromAddressList(addressListData []map[string]string, ip string) {

	for _, list := range addressListData {
		if list["address"] == ip {
			_, err := mr.client.Run(
				"/ip/firewall/address-list/remove",
				fmt.Sprintf("=numbers=%s", list[".id"]),
			)

			if err == nil {
				log.Println(list)
			}
		}
	}

}

func (mr *mikrotikRepository) GetAddressList(listName string) []map[string]string {
	var M []map[string]string
	reply, err := mr.client.Run(
		"/ip/firewall/address-list/print",
	)

	if err != nil {
		log.Println("Error al imprimir los address list", err)
	}

	for _, alist := range reply.Re {
		if alist.Map["list"] == listName {
			M = append(M, alist.Map)
		}
	}
	return M
}
