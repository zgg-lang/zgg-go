package parser

import "github.com/antlr4-go/antlr/v4"

type ZggBaseParser struct {
	*antlr.BaseParser
}

func (p *ZggBaseParser) here(_type int) bool {
	possibleIndexEosToken := p.GetCurrentToken().GetTokenIndex() - 1
	ahead := p.GetTokenStream().Get(possibleIndexEosToken)
	return ahead.GetChannel() == antlr.LexerHidden && ahead.GetTokenType() == _type
}

// func (p *ZggBaseParser) notLineTerminator() bool {
// 	return !p.here(ZggLexerNEWLINE)
// }
