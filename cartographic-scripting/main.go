package main

import "go.starlark.net/starlark"
import "fmt"
import "os"

type Feature struct {
	id int
	properties map[string]string
}

type Layer struct {
	name string
	filter starlark.Callable
	tags *starlark.List
}

type ComplexTag struct {
	key string
	value starlark.Value
}

func (t ComplexTag) Freeze() {

}

func (t ComplexTag) Hash() (uint32,error) {
	return 0,nil
}

func (t ComplexTag) String() string {
	return t.key
}

func (t ComplexTag) Truth() starlark.Bool {
	return true
}

func (t ComplexTag) Type() string {
	return "complextag"
}

func main() {
	thread := &starlark.Thread{Name: "my thread"}

	var layers []Layer

	layerBuiltin := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var layer Layer
		// var name string
		// var filter starlark.Callable
		// var tags *starlark.List
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "name", &layer.name, "tags?", &layer.tags, "filter?", &layer.filter); err != nil {
			return nil, err
		}
		layers = append(layers,layer)
		return starlark.MakeInt(0), nil
	}

	tagBuiltin := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var tag ComplexTag
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "key", &tag.key, "value", &tag.value); err != nil {
			return nil, err
		}
		return tag, nil
	}

	predeclared := starlark.StringDict{
		"layer":   starlark.NewBuiltin("layer", layerBuiltin),
		"tag":   starlark.NewBuiltin("tag", tagBuiltin),
	}

	globals, err := starlark.ExecFile(thread, "tileset_definition.star", nil, predeclared)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	attribution := globals["ATTRIBUTION"]
	fmt.Println(attribution)
	fmt.Println(attribution.Type())
	fmt.Println(layers)

	features := []Feature{
		{
			id: 1,
			properties: map[string]string{
					"building":"yes",
					"height":"20",
			},
		},
		{
			id: 2,
			properties: map[string]string{
					"highway":"motorway",
			},
		},
		{
			id: 3,
			properties: map[string]string{
					"highway":"secondary",
			},
		},
	}

	// for every feature in the input...
	for _, feature := range features {
		feature_dict := starlark.NewDict(1)
		props_dict := starlark.NewDict(1)
		for key, elem := range feature.properties {
			props_dict.SetKey(starlark.String(key),starlark.String(elem))
		}
		feature_dict.SetKey(starlark.String("properties"),props_dict)

		// for every zoom level...
		for zoom := 0; zoom <= 4; zoom++ {

			for _, layer := range layers {
				fmt.Println(layer)
				v, err := starlark.Call(thread, layer.filter, starlark.Tuple{starlark.MakeInt(zoom),feature_dict}, nil)
				if err == nil && v.Truth() {
					fmt.Println("feature", feature.id, "appears in",layer.name,"at zoom", zoom)
				}
			}

			// iter := layers_list.Iterate()
			// defer iter.Done()
			// var layer starlark.Value

			// // for every layer defined...
			// for iter.Next(&layer) {
			// 	layer_dict, _ := layer.(*starlark.Dict)
			// 	name_val, _, _ := layer_dict.Get(starlark.String("name"))
			// 	filter_val, _, _ := layer_dict.Get(starlark.String("filter"))
			// 	filter_func, _ := filter_val.(*starlark.Function)
			// }
		}
	}
}
