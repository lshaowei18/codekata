package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

/*
Schema:
- char  -Boolean arg
- char* -String arg
- char# -Integer arg
*/

type Parser struct {
	flags  []string
	args   []string
	schema map[string]interface{} //Key value pair, value is default value.
}

func Args(flags string, args []string) (Parser, error) {
	split := strings.Split(flags, ",") //Get individual flags
	args = args[1:]
	p := Parser{split, args, map[string]interface{}{}}
	err := p.createSchema()
	if err != nil {
		return Parser{}, errors.Wrap(err, "Error creating Schema")
	}
	err = p.parseArgs()
	if err != nil {
		return Parser{}, errors.Wrap(err, "Error parsing arguments")
	}
	return p, nil
}

func (p *Parser) GetValue(arg string) (interface{}, error) {
	val, ok := p.schema[arg]
	if !ok {
		return nil, fmt.Errorf("Arg is not in schema.")
	}
	return val, nil
}

func (p *Parser) createSchema() error {
	for _, flag := range p.flags {

		key, flagType := parseFlag(flag)

		if flagType == reflect.Invalid {
			return fmt.Errorf("Parsing error, invalid flag type.")
		}
		var val interface{}

		switch flagType {
		case reflect.Int:
			val = 0
		case reflect.Bool:
			val = false
		case reflect.String:
			val = ""
		}

		p.schema[key] = val
	}
	return nil
}

func (p *Parser) parseArgs() error {
	for i := 0; i < len(p.args); i++ {
		next := i + 1

		//If next is more than the total number of arguments, break
		if next >= len(p.args) {
			break
		}

		//If next argument is a flag, move on
		if p.args[next][0] == '-' {
			continue
		}

		key := string(p.args[i][1])

		flagType, err := p.getFlagType(key)
		if err != nil {
			return errors.Wrap(err, "Error parsing args")
		}

		val, err := convertStringToType(p.args[next], flagType)
		if err != nil {
			return errors.Wrap(err, "Can't convert arg to respective type")
		}

		p.schema[key] = val
	}
	return nil
}

func (p Parser) getFlagType(flag string) (reflect.Kind, error) {
	v, ok := p.schema[flag]
	if !ok {
		return reflect.Invalid, fmt.Errorf("Given arg %v not in schema.", flag)
	}
	return reflect.ValueOf(v).Kind(), nil
}

func parseFlag(flag string) (string, reflect.Kind) {
	//Boolean flags
	if len(flag) == 1 {
		return flag, reflect.Bool
	}

	key := string(flag[0])

	if flag[1] == '*' {
		return key, reflect.String
	} else if flag[1] == '#' {
		return key, reflect.Int
	}
	return key, reflect.Invalid
}

func (p Parser) printSchema() {
	for k, v := range p.schema {
		fmt.Printf("Key %v, value %v of type %T\n", k, v, v)
	}
}

func convertStringToType(val string, flag reflect.Kind) (interface{}, error) {
	switch flag {
	case reflect.Int:
		value, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return value, nil
	case reflect.Bool:
		value, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return value, nil
	default: //We assume the default is string
		return val, nil
	}
}
