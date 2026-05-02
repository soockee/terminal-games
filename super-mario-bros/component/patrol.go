package component

import "github.com/yohamta/donburi"

// PatrolData drives an entity along a sequence of waypoints.
// Waypoints are resolved from LDtk Path_node linked lists at load time.
type PatrolData struct {
	Waypoints []Vec2  // world-space points [start, node1, node2, ...]
	Current   int     // index of the waypoint we're heading toward
	Speed     float64 // pixels per second
	Forward   bool    // true = advancing through waypoints, false = reversing
	Loop      bool    // true = jump back to start at end, false = ping-pong
}

// Vec2 is a simple 2D point (avoids importing cp in the component package).
type Vec2 struct {
	X, Y float64
}

var Patrol = donburi.NewComponentType[PatrolData]()
