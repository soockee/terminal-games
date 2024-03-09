package helper

import "github.com/solarlune/ldtkgo"

func GetEntityByName(name string, level int, project *ldtkgo.Project) *ldtkgo.Entity {
	for _, layer := range project.Levels[level].Layers {
		for _, entity := range layer.Entities {
			if entity.Identifier == name {
				return entity
			}
		}
	}
	return nil
}
