package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ManagementClient struct {
	Config *Config
}

func NewManagementClient(cfg *Config) *ManagementClient {
	return &ManagementClient{cfg}
}

func (c *ManagementClient) GetMappings() ([]Mapping, error) {
	addr := fmt.Sprintf("http://%v/mappings/", c.Config.MgmtAddr)

	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mappings := make([]Mapping, 0)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (c *ManagementClient) AddMapping(m *Mapping) error {
	addr := fmt.Sprintf("http://%v/mappings/", c.Config.MgmtAddr)

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	resp, err := http.Post(addr, "application/json", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}

func (c *ManagementClient) RemoveMapping(m *Mapping) error {
	addr := fmt.Sprintf("http://%v/mappings/%v", c.Config.MgmtAddr, m.Key)

	req, err := http.NewRequest("DELETE", addr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
func (c *ManagementClient) RemoveAllMappings() error {
	addr := fmt.Sprintf("http://%v/mappings/", c.Config.MgmtAddr)

	req, err := http.NewRequest("DELETE", addr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}
