package main

import "go.starlark.net/starlark"
import "fmt"
import "os"

type Feature struct {
	id int
	properties map[string]string
}

func main() {
	thread := &starlark.Thread{Name: "my thread"}
	globals, err := starlark.ExecFile(thread, "tileset_definition.star", nil, nil)
	if err != nil {
		fmt.Println("oops")
		os.Exit(1)
	}

	attribution := globals["ATTRIBUTION"]
	fmt.Println(attribution)
	fmt.Println(attribution.Type())

	layers := globals["LAYERS"]
	layers_list, ok := layers.(*starlark.List)
	if !ok {
		fmt.Println("Layers must be list")
		os.Exit(1)
	}

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
			// prepare the call

			iter := layers_list.Iterate()
			defer iter.Done()
			var layer starlark.Value

			// for every layer defined...
			for iter.Next(&layer) {
				layer_dict, _ := layer.(*starlark.Dict)
				name_val, _, _ := layer_dict.Get(starlark.String("name"))
				filter_val, _, _ := layer_dict.Get(starlark.String("filter"))
				filter_func, _ := filter_val.(*starlark.Function)

				v, err := starlark.Call(thread, filter_func, starlark.Tuple{starlark.MakeInt(zoom),feature_dict}, nil)
				if err == nil && v.Truth() {
					fmt.Println("feature", feature.id, "appears in",name_val,"at zoom", zoom)
				}
			}
		}
	}
}
