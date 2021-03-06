package components

import (
	"os"
	"plugin"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/config"
	"github.com/fananchong/go-xserver/internal/components/misc"
)

// Plugin : 插件组件
type Plugin struct {
	ctx       *common.Context
	pluginObj common.IPlugin
}

// NewPlugin : 实例化
func NewPlugin(ctx *common.Context) *Plugin {
	p := &Plugin{ctx: ctx}
	p.loadPlugin(ctx)
	return p
}

// Start : 实例化组件
func (p *Plugin) Start() bool {
	var ret bool
	if p.pluginObj != nil {
		ret = p.pluginObj.Start()
	}
	return ret
}

// Close : 关闭组件
func (p *Plugin) Close() {
	if p.pluginObj != nil {
		p.pluginObj.Close()
		p.pluginObj = nil
	}
}

func (p *Plugin) loadPlugin(ctx *common.Context) {
	v := p.ctx.IConfig.(*Config).GetViperObj()
	appName := v.GetString("app")
	if appName == "" {
		p.ctx.PrintUsage()
		os.Exit(1)
	}
	so, err := plugin.Open(appName + ".so")
	if err != nil {
		ctx.Errorln(err)
		os.Exit(1)
	}
	obj, err := so.Lookup("PluginObj")
	if err != nil {
		ctx.Errorln(err)
		os.Exit(1)
	}
	t, err := so.Lookup("PluginType")
	if err != nil {
		ctx.Errorln(err)
		os.Exit(1)
	}
	c, err := so.Lookup("Ctx")
	if err != nil {
		ctx.Errorln(err)
		os.Exit(1)
	}
	p.pluginObj = *obj.(*common.IPlugin)
	pluginType := *t.(*config.NodeType)
	*c.(**common.Context) = ctx
	misc.SetPluginType(ctx, pluginType)
}
