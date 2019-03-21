package parsers

import (
	"strings"

	"github.com/haproxytech/config-parser/common"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/types"
)

type OptionForwardFor struct {
	data *types.OptionForwardFor
}

/*
option forwardfor [ except <network> ] [ header <name> ] [ if-none ]
*/
func (s *OptionForwardFor) Parse(line string, parts, previousParts []string, comment string) (changeState string, err error) {
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "forwardfor" {
		data := &types.OptionForwardFor{
			Comment: comment,
		}
		index := 2
		for index < len(parts) {
			switch parts[index] {
			case "except":
				index++
				if index == len(parts) {
					return "", errors.InvalidData
				}
				data.Except = parts[index]
			case "header":
				index++
				if index == len(parts) {
					return "", errors.InvalidData
				}
				data.Header = parts[index]
			case "if-none":
				data.IfNone = true
			default:
				return "", errors.InvalidData
			}
			index++
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option forwardfor", Line: line}
}

func (s *OptionForwardFor) Result(AddComments bool) ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.FetchError
	}
	var sb strings.Builder
	sb.WriteString("option forwardfor")
	//option forwardfor [ except <network> ] [ header <name> ] [ if-none ]
	if s.data.Except != "" {
		sb.WriteString(" except ")
		sb.WriteString(s.data.Except)
	}
	if s.data.Header != "" {
		sb.WriteString(" header ")
		sb.WriteString(s.data.Header)
	}
	if s.data.IfNone {
		sb.WriteString(" if-none")
	}
	return []common.ReturnResultLine{
		common.ReturnResultLine{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}