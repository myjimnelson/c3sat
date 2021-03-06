package queryciv3

import (
	"encoding/hex"

	"github.com/graphql-go/graphql"
)

const tileBytes = 212

type worldData struct {
	worldOffset int
}

type mapData struct {
	mapWidth          int
	mapHeight         int
	tileSetWidth      int
	tileSetHeight     int
	tileSetX          int
	tileSetY          int
	playerSpoilerMask int32
	tilesOffset       int
	tileSetOffsets    []int
	mapTileOffsets    []int
}

func (m *mapData) spoilerFree(offset int) bool {
	if m.playerSpoilerMask == 0 || int(m.playerSpoilerMask)&ReadInt32(offset+84, Unsigned) != 0 {
		return true
	}
	return false
}

// To return to subqueries in GraphQL
type saveAndOffsetType struct {
	save    *saveGameType
	offset  int
	offset2 int
}

var gameLeadSectionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "gameLeadSection",
	Fields: graphql.Fields{
		"bytes": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Byte array",
			Args: graphql.FieldConfigArgument{
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of item",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int32s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if itemOffset, ok := p.Source.(int); ok {
					if itemOffset > 0 {
						offset, _ := p.Args["offset"].(int)
						count, _ := p.Args["count"].(int)
						intList := make([]int, count)
						for i := 0; i < count; i++ {
							intList[i] = ReadInt8((itemOffset+4+offset)+i, Unsigned)
						}
						return intList, nil
					}
				}
				return nil, nil
			},
		},
		"int32s": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Int32 array",
			Args: graphql.FieldConfigArgument{
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of item",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int32s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if itemOffset, ok := p.Source.(int); ok {
					if itemOffset > 0 {
						offset, _ := p.Args["offset"].(int)
						count, _ := p.Args["count"].(int)
						intList := make([]int, count)
						for i := 0; i < count; i++ {
							intList[i] = ReadInt32((itemOffset+4+offset)+4*i, Signed)
						}
						return intList, nil
					}
				}
				return nil, nil
			},
		},
		"hexDump": &graphql.Field{
			Type:        graphql.String,
			Description: "Hex dump of the entire item",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if itemOffset, ok := p.Source.(int); ok {
					if itemOffset > 0 {
						length := ReadInt32(itemOffset, Signed)
						return hex.Dump(saveGame.data[itemOffset+4 : itemOffset+4+length]), nil
					}
				}
				return "", nil
			},
		},
	},
})

// intending to append fields for city lists and mgl lists, but not in a hurry
var raceSectionItemType = listSectionItem

var listSectionItem = graphql.NewObject(graphql.ObjectConfig{
	Name: "listSectionItem",
	Fields: graphql.Fields{
		"bytes": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Byte array",
			Args: graphql.FieldConfigArgument{
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of item",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int32s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(saveAndOffsetType); ok {
					if item.offset > 0 {
						offset, _ := p.Args["offset"].(int)
						count, _ := p.Args["count"].(int)
						intList := make([]int, count)
						for i := 0; i < count; i++ {
							intList[i] = item.save.readInt8((item.offset+4+offset)+4*i, Unsigned)
						}
						return intList, nil
					}
				}
				return nil, nil
			},
		},
		"int32s": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Int32 array",
			Args: graphql.FieldConfigArgument{
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of item",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int32s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(saveAndOffsetType); ok {
					if item.offset > 0 {
						offset, _ := p.Args["offset"].(int)
						count, _ := p.Args["count"].(int)
						intList := make([]int, count)
						for i := 0; i < count; i++ {
							intList[i] = item.save.readInt32((item.offset+4+offset)+4*i, Signed)
						}
						return intList, nil
					}
				}
				return nil, nil
			},
		},
		"hexDump": &graphql.Field{
			Type:        graphql.String,
			Description: "Hex dump of the entire item",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(saveAndOffsetType); ok {
					if item.offset > 0 {
						length := item.save.readInt32(item.offset, Signed)
						return hex.Dump(saveGame.data[item.offset+4 : item.offset+4+length]), nil
					}
				}
				return "", nil
			},
		},
		"string": &graphql.Field{
			Type:        graphql.String,
			Description: "Null-terminated string",
			Args: graphql.FieldConfigArgument{
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of item",
				},
				"maxLength": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Max length of string / the max number of bytes to consider",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if item, ok := p.Source.(saveAndOffsetType); ok {
					if item.offset > 0 {
						offset, _ := p.Args["offset"].(int)
						maxLength, _ := p.Args["maxLength"].(int)
						stringBuffer := item.save.data[item.offset+4+offset : item.offset+4+offset+maxLength]
						s, err := CivString(stringBuffer)
						if err != nil {
							return nil, err
						}
						return s, nil
					}
				}
				return "", nil
			},
		},
	},
})

var mapTileType = graphql.NewObject(graphql.ObjectConfig{
	Name: "tile",
	Fields: graphql.Fields{
		"hexTerrain": &graphql.Field{
			Type:        graphql.String,
			Description: "Byte value. High nybble is overlay, low nybble is base terrain",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if offset, ok := p.Source.(int); ok {
					if offset > 0 {
						return hex.EncodeToString(saveGame.data[offset+57 : offset+58]), nil
					}
				}
				return nil, nil
			},
		},
		"chopped": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "True if a forest has previously been harvested from this tile",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if offset, ok := p.Source.(int); ok {
					if offset > 0 {
						return ((ReadInt16(offset+62, Unsigned) & 0x1000) != 0), nil
					}
				}
				return nil, nil
			},
		},
	},
})

var mapType = graphql.NewObject(graphql.ObjectConfig{
	Name: "map",
	Fields: graphql.Fields{
		"mapWidth": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width of the game map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.mapWidth, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"mapHeight": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height of the game map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.mapHeight, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetWidth": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width of the currently visible map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetWidth, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetHeight": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height of the currently visible map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetHeight, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetX": &graphql.Field{
			Type:        graphql.Int,
			Description: "World map X coordinate of top-left tile set tile",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetX, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetY": &graphql.Field{
			Type:        graphql.Int,
			Description: "World map Y coordinate of top-left tile set tile",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetY, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tiles": &graphql.Field{
			Type:        graphql.NewList(mapTileType),
			Description: "List of all visible tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetOffsets, nil
				}
				return nil, nil
			},
		},
	},
})

var civ3Type = graphql.NewObject(graphql.ObjectConfig{
	Name: "civ3",
	Fields: graphql.Fields{
		"worldSeed": &graphql.Field{
			Type:        graphql.Int,
			Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+170, Signed), nil
				}
				return -1, nil
			},
		},
		"climate": &graphql.Field{
			Type:        graphql.Int,
			Description: "Climate option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+186, Signed), nil
				}
				return -1, nil
			},
		},
		"climateFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Climate setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+190, Signed), nil
				}
				return -1, nil
			},
		},
		"barbarians": &graphql.Field{
			Type:        graphql.Int,
			Description: "Barbarians option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+194, Signed), nil
				}
				return -1, nil
			},
		},
		"barbariansFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Barbarians setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+198, Signed), nil
				}
				return -1, nil
			},
		},
		"landMass": &graphql.Field{
			Type:        graphql.Int,
			Description: "Land mass option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+202, Signed), nil
				}
				return -1, nil
			},
		},
		"landMassFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Land mass setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+206, Signed), nil
				}
				return -1, nil
			},
		},
		"oceanCoverage": &graphql.Field{
			Type:        graphql.Int,
			Description: "Ocean coverage option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+210, Signed), nil
				}
				return -1, nil
			},
		},
		"oceanCoverageFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Ocean coverage setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+214, Signed), nil
				}
				return -1, nil
			},
		},
		"temperature": &graphql.Field{
			Type:        graphql.Int,
			Description: "Temperature option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+218, Signed), nil
				}
				return -1, nil
			},
		},
		"temperatureFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Temperature setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+222, Signed), nil
				}
				return -1, nil
			},
		},
		"age": &graphql.Field{
			Type:        graphql.Int,
			Description: "Age option chosen for random map generation",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+226, Signed), nil
				}
				return -1, nil
			},
		},
		"ageFinal": &graphql.Field{
			Type:        graphql.Int,
			Description: "Age setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+230, Signed), nil
				}
				return -1, nil
			},
		},
		"size": &graphql.Field{
			Type:        graphql.Int,
			Description: "Size setting of map",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+234, Signed), nil
				}
				return -1, nil
			},
		},
	},
})
