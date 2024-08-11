package lilv

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/libs -l:liblilv-0.so

#include <stdlib.h>
#include <lv2/atom/atom.h>
#include <lv2/core/lv2.h>
#include <lv2/urid/urid.h>
#include <lv2/worker/worker.h>
#include <lilv/lilv.h>

extern int go_urid_map(void*, char*);

extern int go_schedule_work(void*, uint32_t, void*);

static inline LV2_URID map_uri(LV2_URID_Map_Handle handle, const char * uri) {
	const LV2_URID id = go_urid_map(handle, (char*)uri);
	return id;
}

static inline LV2_URID_Map* new_urid_map(LV2_URID_Map_Handle handle) {
	LV2_URID_Map* self = malloc(sizeof(LV2_URID_Map));
	self->handle = handle;
	self->map =  map_uri;
	return self;
}

static inline LV2_Worker_Status schedule_work(LV2_Worker_Schedule_Handle handle, uint32_t size, const void* data) {
	const LV2_Worker_Status status = go_schedule_work(handle, size, (void*)data);
	return status;
}

static inline LV2_Worker_Schedule* new_worker_schedule(LV2_Worker_Schedule_Handle handle) {
	LV2_Worker_Schedule* self = malloc(sizeof(LV2_Worker_Schedule));
	self->handle = handle;
	self->schedule_work = schedule_work;
	return self;
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	_ "github.com/ianlancetaylor/cgosymbolizer"
)

type World struct {
	world *C.LilvWorld
}

func NewWorld() *World {
	return &World{
		world: C.lilv_world_new(),
	}
}

func (w *World) Free() {
	C.lilv_world_free(w.world)
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

func (p *Plugin) Instantiate(sampleRate float64, features []LV2Feature) *Instance {
	if p.plugin == nil {
		return nil
	}

	lv2Features := make([]*C.LV2_Feature, 0)

	pins := make([]runtime.Pinner, 0)
	for i, f := range features {
		var p runtime.Pinner
		pins = append(pins, p)
		var feature C.LV2_Feature
		feature.URI = C.CString(f.URI)
		feature.data = f.Data()
		lv2Features = append(lv2Features, &feature)
		p.Pin(lv2Features[i])
	}
	defer func() {
		for _, p := range pins {
			p.Unpin()
		}
	}()

	instance := C.lilv_plugin_instantiate(p.plugin, (C.double)(sampleRate), (**C.LV2_Feature)(unsafe.Pointer(unsafe.SliceData(lv2Features))))
	if instance == nil {
		return nil
	}
	return &Instance{
		instance: instance,
	}
}

type Instance struct {
	instance *C.LilvInstance
}

func (i *Instance) ConnectPort(index int, data unsafe.Pointer) {
	C.lilv_instance_connect_port(i.instance, C.uint(index), data)
}

func (i *Instance) Activate() {
	C.lilv_instance_activate(i.instance)
}

func (i *Instance) Run(sampleCount uint) {
	C.lilv_instance_run(i.instance, (C.uint)(sampleCount))
}

func (i *Instance) Deactivate() {
	C.lilv_instance_deactivate(i.instance)
}

func (i *Instance) Free() {
	C.lilv_instance_free((i.instance))
}

type LV2Feature struct {
	URI  string
	data string
}

func NewLV2Feature(URI string, dataJSON string) LV2Feature {
	return LV2Feature{
		URI:  URI,
		data: dataJSON,
	}
}

func (f LV2Feature) Data() unsafe.Pointer {
	switch f.URI {
	case "http://lv2plug.in/ns/ext/urid#map":
		return unsafe.Pointer(C.new_urid_map((C.LV2_URID_Map_Handle)(unsafe.Pointer(&f))))
	case "http://lv2plug.in/ns/ext/worker#schedule":
		return unsafe.Pointer(C.new_worker_schedule((C.LV2_Worker_Schedule_Handle)(unsafe.Pointer(&f))))
	}
	return nil
}

var uridMap = make(map[string]uint32)

//export go_urid_map
func go_urid_map(p unsafe.Pointer, p1 *C.char) C.int {
	//feature := (*LV2Feature)(p)
	str := C.GoString(p1)
	var id uint32
	var ok bool
	if id, ok = uridMap[str]; !ok {
		id = uint32(len(uridMap)) + 1
		uridMap[str] = id
	}
	fmt.Printf("go_map_urid: %s - id: %d\n", str, id)
	return (C.int)(id)
}

//export go_schedule_work
func go_schedule_work(p unsafe.Pointer, size uint32, data unsafe.Pointer) C.int {
	//feature := (*LV2Feature)(p)
	fmt.Printf("go_schedule_work: size: %v\n", size)
	return 0
}
