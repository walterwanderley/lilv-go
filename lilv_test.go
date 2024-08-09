package lilv_test

import (
	"testing"

	"github.com/walterwanderley/lilv-go"
)

func TestNewWorld(t *testing.T) {
	w := lilv.NewWorld()
	t.Logf("World: %T %+v", w, w)

	w.SetLv2Path("./testdata")

	w.LoadBundle("testdata/neural_amp_modeler.lv2")

	w.LoadSpecifications()
	w.LoadPluginClasses()

	nodes := w.FindNodes("http://github.com/mikeoliphant/neural-amp-modeler-lv2", "http://usefulinc.com/ns/doap#name", "")

	t.Logf("Nodes: %T %+v", nodes, nodes)

	node := w.Get("http://github.com/mikeoliphant/neural-amp-modeler-lv2", "http://usefulinc.com/ns/doap#name", "")

	t.Logf("Node: %T %+v", node, node)

	ask := w.Ask("http://github.com/mikeoliphant/neural-amp-modeler-lv2", "http://usefulinc.com/ns/doap#name", "")
	t.Logf("Ask, %v", ask)

	plugins := w.GetAllPlugins()
	t.Logf("plugins: %v, %d", plugins, plugins.Size())

	plugin := plugins.GetByURI("http://github.com/mikeoliphant/neural-amp-modeler-lv2")
	t.Logf("plugin: %v", plugin)

	instance := plugin.Instantiate(48000.0, []*lilv.LV2Feature{
		{
			URI: "http://lv2plug.in/ns/ext/urid#map",
		},
	})

	t.Logf("Intance: %#v", instance)
}
