package engine

import "github.com/solarlune/resolv"

func IsInside(test resolv.IShape, in resolv.IShape) bool {
	switch s := test.(type) {
	case *resolv.ConvexPolygon:
		for _, p := range s.Points {
			if !p.IsInside(in) {
				return false
			}
		}
		return true
	case *resolv.Circle:
		return s.Position().IsInside(in)
	}
	return false
}
