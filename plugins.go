package gokong

import (
	"encoding/json"
	"fmt"
)

type PluginClient interface {
	GetById(id string) (*Plugin, error)
	List(query *PluginQueryString) ([]*Plugin, error)
	Create(pluginRequest *PluginRequest) (*Plugin, error)
	UpdateById(id string, pluginRequest *PluginRequest) (*Plugin, error)
	DeleteById(id string) error
	GetByConsumerId(id string) (*Plugins, error)
	GetByRouteId(id string) (*Plugins, error)
	GetByServiceId(id string) (*Plugins, error)
}

type pluginClient struct {
	config *Config
}

type PluginRequest struct {
	Name       string                 `json:"name" yaml:"name"`
	ConsumerId *Id                    `json:"consumer" yaml:"consumer"`
	ServiceId  *Id                    `json:"service" yaml:"service"`
	RouteId    *Id                    `json:"route" yaml:"route"`
	RunOn      string                 `json:"run_on,omitempty" yaml:"run_on,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type Plugin struct {
	Id         string                 `json:"id" yaml:"id"`
	Name       string                 `json:"name" yaml:"name"`
	ConsumerId *Id                    `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	ServiceId  *Id                    `json:"service,omitempty" yaml:"service,omitempty"`
	RouteId    *Id                    `json:"route,omitempty" yaml:"route,omitempty"`
	RunOn      string                 `json:"run_on,omitempty" yaml:"run_on,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type Plugins struct {
	Data   []*Plugin `json:"data" yaml:"data,omitempty"`
	Next   *string   `json:"next" yaml:"next,omitempty"`
	Offset string    `json:"offset,omitempty" yaml:"offset,omitempty"`
}

type PluginQueryString struct {
	Offset string `json:"offset,omitempty" yaml:"offset,omitempty"`
	Size   int    `json:"size" yaml:"size,omitempty"`
}

const PluginsPath = "/plugins/"

func (pluginClient *pluginClient) GetById(id string) (*Plugin, error) {

	r, body, errs := newGet(pluginClient.config, PluginsPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugin := &Plugin{}
	err := json.Unmarshal([]byte(body), plugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin plugin response, error: %v", err)
	}

	if plugin.Id == "" {
		return nil, nil
	}

	return plugin, nil
}

func (pluginClient *pluginClient) List(query *PluginQueryString) ([]*Plugin, error) {
	plugins := make([]*Plugin, 0)

	if query.Size < 100 {
		query.Size = 100
	}

	if query.Size > 1000 {
		query.Size = 1000
	}

	for {
		data := &Plugins{}

		r, body, errs := newGet(pluginClient.config, PluginsPath).Query(*query).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get plugins, error: %v", errs)
		}

		if r.StatusCode == 401 || r.StatusCode == 403 {
			return nil, fmt.Errorf("not authorised, message from kong: %s", body)
		}

		err := json.Unmarshal([]byte(body), data)
		if err != nil {
			return nil, fmt.Errorf("could not parse plugins list response, error: %v", err)
		}

		plugins = append(plugins, data.Data...)

		if data.Next == nil || *data.Next == "" {
			break
		}

		query.Offset = data.Offset
	}

	return plugins, nil
}

func (pluginClient *pluginClient) Create(pluginRequest *PluginRequest) (*Plugin, error) {
	r, body, errs := newPost(pluginClient.config, PluginsPath).Send(pluginRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	if r.StatusCode == 400 {
		return nil, fmt.Errorf("bad request, message from kong: %s", body)
	}

	createdPlugin := &Plugin{}
	err := json.Unmarshal([]byte(body), createdPlugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin creation response, error: %v kong response: %s", err, body)
	}

	if createdPlugin.Id == "" {
		return nil, fmt.Errorf("could not create plugin, err: %v", body)
	}

	return createdPlugin, nil
}

func (pluginClient *pluginClient) UpdateById(id string, pluginRequest *PluginRequest) (*Plugin, error) {
	r, body, errs := newPatch(pluginClient.config, PluginsPath+id).Send(pluginRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	if r.StatusCode == 400 {
		return nil, fmt.Errorf("bad request, message from kong: %s", body)
	}

	updatedPlugin := &Plugin{}
	err := json.Unmarshal([]byte(body), updatedPlugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin update response, error: %v kong response: %s", err, body)
	}

	if updatedPlugin.Id == "" {
		return nil, fmt.Errorf("could not update plugin, error: %v", body)
	}

	return updatedPlugin, nil
}

func (pluginClient *pluginClient) DeleteById(id string) error {
	r, body, errs := newDelete(pluginClient.config, PluginsPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete plugin, result: %v error: %v", r, errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return fmt.Errorf("not authorised, message from kong: %s", body)
	}

	return nil
}

func (pluginClient *pluginClient) GetByConsumerId(id string) (*Plugins, error) {
	r, body, errs := newGet(pluginClient.config, "/consumers/"+id+"/plugins").End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugins, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugins := &Plugins{}
	err := json.Unmarshal([]byte(body), plugins)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugins list response, error: %v", err)
	}

	return plugins, nil
}

func (pluginClient *pluginClient) GetByRouteId(id string) (*Plugins, error) {
	r, body, errs := newGet(pluginClient.config, "/routes/"+id+"/plugins").End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugins, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugins := &Plugins{}
	err := json.Unmarshal([]byte(body), plugins)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugins list response, error: %v", err)
	}

	return plugins, nil
}

func (pluginClient *pluginClient) GetByServiceId(id string) (*Plugins, error) {
	r, body, errs := newGet(pluginClient.config, "/services/"+id+"/plugins").End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugins, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugins := &Plugins{}
	err := json.Unmarshal([]byte(body), plugins)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugins list response, error: %v", err)
	}

	return plugins, nil
}
