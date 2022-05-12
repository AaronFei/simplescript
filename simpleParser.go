package simplescript

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ScriptRunner_t struct {
	translateTable reflect.Value
	funcSet        reflect.Value
	showlog        bool
}

func CreateScriptRunner() ScriptRunner_t {
	return ScriptRunner_t{
		translateTable: reflect.ValueOf(nil),
		funcSet:        reflect.ValueOf(nil),
		showlog:        false,
	}
}

func (s *ScriptRunner_t) RegisterFuncSet(funcSet interface{}) {
	funcSet = reflect.ValueOf(funcSet)
}

func (s *ScriptRunner_t) InstallFunc(name string, funcName string) {
	s.translateTable.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(funcName))
}

func (s *ScriptRunner_t) Uninstall() {
	s.translateTable = reflect.ValueOf(nil)
	s.funcSet = reflect.ValueOf(nil)
}

func (s *ScriptRunner_t) EnableLog() {
	s.showlog = true
}

func (s *ScriptRunner_t) DisableLog() {
	s.showlog = false
}

func (s *ScriptRunner_t) Execute(commands []string) error {

	if !s.funcSet.IsValid() || !s.translateTable.IsValid() {
		return errors.New("Install table and function set first")
	}

	for _, str := range commands {

		ele := strings.Split(str, " ")

		for i, v := range ele {
			ele[i] = strings.TrimSpace(v)
		}

		funcName := s.translateTable.MapIndex(reflect.ValueOf(ele[0]))
		if !funcName.IsValid() {
			return errors.New("Invalid key " + ele[0])
		}

		pfunc := s.funcSet.MethodByName(funcName.Interface().(string))

		if !pfunc.IsValid() {
			return errors.New("Invalid value " + funcName.Interface().(string))
		}

		para := []reflect.Value{}
		for _, v := range ele[1:] {
			para = append(para, reflect.ValueOf(v))
		}

		if s.showlog {
			fmt.Println(ele[0], para)
		}

		if result := pfunc.Call(para); !result[0].IsNil() {
			return fmt.Errorf("[%s] Error: %v", funcName.Interface().(string), result[0].Interface().(error))
		}
	}

	return nil
}
