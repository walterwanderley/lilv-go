package lilv_test

import (
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/go-audio/wav"
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

	instance := plugin.Instantiate(44100.0, []lilv.LV2Feature{
		lilv.NewLV2Feature("http://lv2plug.in/ns/ext/urid#map", ``),
		lilv.NewLV2Feature("http://lv2plug.in/ns/ext/worker#schedule", `bhsasb`),
	})

	t.Logf("Intance: %#v", instance)
	/*
		inputAudio, err := os.ReadFile("testdata/test.wav")
		if err != nil {
			t.Error(err)

			return
		}
	*/
	f, err := os.Open("/tmp/sinco.wav")
	if err != nil {
		t.Error(err)
		return
	}
	dec := wav.NewDecoder(f)
	dec.ReadMetadata()
	t.Logf("metadata: %+v", dec.Metadata)

	inNode := w.NewFileURI("testdata/test.wav")

	//TODO aprender como configurar esses controls
	//TODO aprender como chamar o Set ou Patch do LV2 para configurar o model nam patch:writable <@NAM_LV2_ID@#model>;
	control := 56.0
	notify := 0.0
	outputAudio := make([]byte, 512)
	inputLevel := 10.0
	outputLevel := 4.0
	instance.PatchSet("http://github.com/mikeoliphant/neural-amp-modeler-lv2#model", "testdata/test.nam")
	instance.ConnectPort(0, unsafe.Pointer(&control))
	instance.ConnectPort(1, unsafe.Pointer(&notify))
	instance.ConnectPort(2, unsafe.Pointer(inNode.Get()))
	instance.ConnectPort(3, unsafe.Pointer(&outputAudio[0]))
	instance.ConnectPort(4, unsafe.Pointer(&inputLevel))
	instance.ConnectPort(5, unsafe.Pointer(&outputLevel))
	instance.Activate()
	instance.Run(512)
	fmt.Printf("out %v", outputAudio)
	instance.Deactivate()
	instance.Free()
	w.Free()

}
