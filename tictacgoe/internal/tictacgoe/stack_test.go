package tictacgoe

import "testing"

func TestPush(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	got := stack.S[0]
	if got != 1 {
		t.Errorf("stack.Push(1) = %d; want 1", got)
	}
}

func TestPushMultiple(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	got := stack.S[1]
	if got != 2 {
		t.Errorf("stack.Push(2) = %d; want 1", got)
	}
}
func TestPop(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	got, err := stack.Pop()
	if err != nil {
		t.Errorf("stack.Pop err %v", err)
	}
	if got != 1 {
		t.Errorf("stack.Pop = %d; want 1", got)
	}
}

func TestPopEmpty(t *testing.T) {
	stack := NewStack[int]()
	_, err := stack.Pop()
	if err == nil {
		t.Error("stack.Pop on empty stack needs to return err")
	}
}

func TestPeek(t *testing.T) {
	stack := NewStack[int]()
	stack.Push(1)
	got, err := stack.Peek()
	if err != nil {
		t.Errorf("stack.Peek err %v", err)
	}
	if got != 1 {
		t.Errorf("stack.Peek = %d; want 1", got)
	}
}

func TestPeekEmpty(t *testing.T) {
	stack := NewStack[int]()
	_, err := stack.Peek()
	if err == nil {
		t.Error("stack.Peek on empty stack needs to return err")
	}
}
