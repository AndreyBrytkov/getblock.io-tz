package utils

import "fmt"

// Sets
type Set struct {
	elements map[string]ElementData
}

type ElementData struct{}

func (s *Set) Add(elemName string, data ElementData) {

	if len(s.elements) == 0 {
		s.elements = make(map[string]ElementData)
	}

	if _, ok := s.elements[elemName]; !ok {
		s.elements[elemName] = data
	}
}

func (s *Set) Del(elemName string) {
	delete(s.elements, elemName)
}

func (s Set) Copy() Set {
	newSet := Set{}
	for elem, data := range s.elements {
		newSet.Add(elem, data)
	}
	return newSet
}

func (s Set) Contains(elemName string) bool {
	if _, ok := s.elements[elemName]; ok {
		return true
	} else {
		return false
	}
}

func (s Set) Union(setB Set) Set {
	newSet := s.Copy()
	for elem, data := range setB.elements {
		newSet.Add(elem, data)
	}
	return newSet
}

func (s Set) Intersection(setB Set) Set {
	var newSet Set
	setA := s
	for elem, data := range setA.elements {
		if setB.Contains(elem) {
			newSet.Add(elem, data)
		}
	}
	return newSet
}

func (s Set) Difference(setB Set) Set {
	newSet := s.Copy()
	for elem := range newSet.elements {
		if setB.Contains(elem) {
			newSet.Del(elem)
		}
	}
	return newSet
}

func (s Set) String() string {
	var str string
	for elem, data := range s.elements {
		str += fmt.Sprintf("%v:%v, ", elem, data)
	}
	return "set(" + str[:len(str)-2] + ")"
}
