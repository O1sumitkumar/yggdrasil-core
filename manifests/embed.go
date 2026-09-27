package manifests

import _ "embed"

//go:embed models/catalog.json
var CatalogJSON []byte

//go:embed models/presets.json
var PresetsJSON []byte
