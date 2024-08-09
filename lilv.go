package lilv

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/libs -l:liblilv-0.so

#include <lv2/core/lv2.h>
#include <lv2/urid/urid.h>
#include <lilv/lilv.h>
*/
import "C"
import (
	"unsafe"
)

type World struct {
	world *C.LilvWorld
}

func NewWorld() *World {
	return &World{
		world: C.lilv_world_new(),
	}
}

func (w *World) SetLv2Path(path string) {
	lv2Path := C.lilv_new_file_uri(w.world, nil, C.CString(path))
	C.lilv_world_set_option(w.world, C.CString(C.LILV_OPTION_LV2_PATH), lv2Path)
}

func (w *World) LoadAll() {
	C.lilv_world_load_all(w.world)
}

func (w *World) LoadBundle(bundleURI string) {
	uri := C.lilv_new_file_uri(w.world, nil, C.CString(bundleURI))
	C.lilv_world_load_bundle(w.world, uri)
}

func (w *World) LoadSpecifications() {
	C.lilv_world_load_specifications(w.world)
}

func (w *World) LoadPluginClasses() {
	C.lilv_world_load_plugin_classes(w.world)
}

func (w *World) NewURI(uri string) *C.LilvNode {
	return C.lilv_new_uri(w.world, C.CString(uri))
}

func (w *World) FindNodes(subject string, predicate string, object string) *Nodes {
	var sub, pre, obj *C.LilvNode
	if subject != "" {
		sub = w.NewURI(subject)
	}
	pre = w.NewURI(predicate)
	if object != "" {
		obj = w.NewURI(object)
	}

	list := C.lilv_world_find_nodes(w.world, sub, pre, obj)
	if list == nil {
		return nil
	}
	return &Nodes{
		nodes: (*C.LilvNodes)(unsafe.Pointer(list)),
	}
}

func (w *World) Get(subject string, predicate string, object string) *Node {
	var sub, pre, obj *C.LilvNode
	if subject != "" {
		sub = w.NewURI(subject)
	}
	pre = w.NewURI(predicate)
	if object != "" {
		obj = w.NewURI(object)
	}
	node := C.lilv_world_get(w.world, sub, pre, obj)
	if node == nil {
		return nil
	}

	return &Node{
		node: (*C.LilvNode)(unsafe.Pointer(node)),
	}
}

func (w *World) Ask(subject string, predicate string, object string) bool {
	var sub, pre, obj *C.LilvNode
	if subject != "" {
		sub = w.NewURI(subject)
	}
	pre = w.NewURI(predicate)
	if object != "" {
		obj = w.NewURI(object)
	}
	result := C.lilv_world_ask(w.world, sub, pre, obj)
	return bool(result)
}

func (w *World) GetAllPlugins() *Plugins {
	plugins := C.lilv_world_get_all_plugins(w.world)
	if plugins == nil {
		return nil
	}
	return &Plugins{
		plugins: (*C.LilvPlugins)(unsafe.Pointer(plugins)),
		world:   w,
	}

}

type Node struct {
	node *C.LilvNode
}

type Nodes struct {
	nodes *C.LilvNodes
}

type Plugins struct {
	plugins *C.LilvPlugins
	world   *World
}

func (p *Plugins) Size() uint {
	if p.plugins == nil {
		return 0
	}
	size := C.lilv_plugins_size(unsafe.Pointer(p.plugins))
	return uint(size)
}

func (p *Plugins) GetByURI(uri string) *Plugin {
	if p.plugins == nil {
		return nil
	}
	plugin := C.lilv_plugins_get_by_uri(unsafe.Pointer(p.plugins), p.world.NewURI(uri))
	if plugin == nil {
		return nil
	}
	return &Plugin{
		plugin: (*C.LilvPlugin)(unsafe.Pointer(plugin)),
	}
}

type Plugin struct {
	plugin *C.LilvPlugin
}

func (p *Plugin) Instantiate(sampleRate float64, features []*LV2Feature) *Instance {
	if p.plugin == nil {
		return nil
	}
	//FIXME add features
	instance := C.lilv_plugin_instantiate((*C.LilvPlugin)(unsafe.Pointer(p.plugin)), (C.double)(sampleRate), nil)
	if instance == nil {
		return nil
	}
	return &Instance{
		instance: (*C.LilvInstance)(unsafe.Pointer(instance)),
	}
}

type Instance struct {
	instance *C.LilvInstance
}

func (i *Instance) ConnectPort(index int, data unsafe.Pointer) {
	C.lilv_instance_connect_port((*C.LilvInstance)(unsafe.Pointer(i.instance)), C.uint(index), data)
}

func (i *Instance) Activate() {
	C.lilv_instance_activate((*C.LilvInstance)(unsafe.Pointer(i.instance)))
}

func (i *Instance) Run(sampleCount uint) {
	C.lilv_instance_run((*C.LilvInstance)(unsafe.Pointer(i.instance)), (C.uint)(sampleCount))
}

func (i *Instance) Deactivate() {
	C.lilv_instance_deactivate((*C.LilvInstance)(unsafe.Pointer(i.instance)))
}

func (i *Instance) Free() {
	C.lilv_instance_free((*C.LilvInstance)(unsafe.Pointer(i.instance)))
}

type LV2Feature struct {
	URI  string
	Data unsafe.Pointer
}
