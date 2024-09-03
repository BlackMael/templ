package parser

import (
	"fmt"
	"os"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goexpression"
)

var conditionalAttribute parse.Parser[ConditionalAttribute] = conditionalAttributeParser{}

type conditionalAttributeParser struct{}

func (conditionalAttributeParser) Parse(pi *parse.Input) (r ConditionalAttribute, ok bool, err error) {
	start := pi.Index()

	// Strip leading whitespace and look for `if `.
	if _, _, err = parse.OptionalWhitespace.Parse(pi); err != nil {
		return
	}
	if !peekPrefix(pi, "if ") {
		pi.Seek(start)
		return
	}

	p, _ := pi.Peek(15)

	// Parse the Go if expression.
	if r.Expression, err = parseGo("if attribute", pi, goexpression.If); err != nil {
		return
	}

	// Eat " {\n".
	if _, ok, err = openBraceWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("attribute if: unterminated (missing closing '{\n')", pi.PositionAt(start))
		return
	}
	if _, _, err = parse.OptionalWhitespace.Parse(pi); err != nil {
		return
	}

	// Read the 'Then' attributes.
	// If there's no match, there's a problem reading the attributes.
	if r.Then, ok, err = (attributesParser{}).Parse(pi); err != nil || !ok {
		err = parse.Error("attribute if: expected attributes in block, but none were found", pi.Position())
		return
	}

	if len(r.Then) == 0 {
		err = parse.Error("attribute if: invalid content or no attributes were found in the if block", pi.Position())
		return
	}

	// Read the optional 'ElseIf' Nodes.
	if r.ElseIfs, _, err = parse.ZeroOrMore(attributeElseIfExpression).Parse(pi); err != nil {
		return
	}

	if len(r.ElseIfs) > 0 {
		fmt.Printf("\n*** We have Else Ifs: %d\n", len(r.ElseIfs))

		for _, elseIf := range r.ElseIfs {
			fmt.Printf("  > %s\n", elseIf.Expression.Value)

			for _, attr := range elseIf.Then {
				// write attr to std out
				_ = attr.Write(os.Stdout, 1)
			}

			fmt.Println("\n        ======== == = =  =")
		}
	}

	// Read the optional 'Else' Nodes.
	if r.Else, ok, err = attributeElseExpression.Parse(pi); err != nil {
		return
	}
	if ok && len(r.Else) == 0 {
		err = parse.Error("attribute if: invalid content or no attributes were found in the else block", pi.Position())
		return
	}

	// Clear any optional whitespace.
	_, _, _ = parse.OptionalWhitespace.Parse(pi)

	fmt.Printf("\nPEEK: \"%s\"", p)
	p, _ = pi.Peek(15)
	fmt.Printf("\n  => \"%s\"\n", p)

	// Read the required closing brace.
	if _, ok, err = closeBraceWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("attribute if: missing end (expected '}')", pi.Position())
		return
	}

	return r, true, nil
}

var attributeElseIfExpression parse.Parser[ElseIfAttribute] = attributeElseIfExpressionParser{}

type attributeElseIfExpressionParser struct{}

func (attributeElseIfExpressionParser) Parse(in *parse.Input) (r ElseIfAttribute, ok bool, err error) {
	start := in.Index()

	// Strip leading whitespace and look for `if `.
	if _, _, err = parse.OptionalWhitespace.Parse(in); err != nil {
		return
	}

	// } else
	var endIfElseParser = parse.All(
		parse.Rune('}'),
		parse.OptionalWhitespace,
		parse.String("else"),
		parse.OptionalWhitespace)
	if _, ok, err = endIfElseParser.Parse(in); err != nil || !ok {
		in.Seek(start)
		return
	}

	if !peekPrefix(in, "if ") {
		in.Seek(start)
		return r, false, nil
	}

	// Parse the Go if expresion.
	if r.Expression, err = parseGo("(else) if attribute", in, goexpression.If); err != nil {
		fmt.Println("Failed to parse Go expression")
		return r, false, err
	}

	// Eat " {\n".
	if _, ok, err = parse.All(openBraceWithOptionalPadding, parse.NewLine).Parse(in); err != nil || !ok {
		err = parse.Error("attribute else if: "+unterminatedMissingCurly, in.PositionAt(start))
		return
	}
	if _, _, err = parse.OptionalWhitespace.Parse(in); err != nil {
		return
	}

	// Read the 'Then' attributes.
	// If there's no match, there's a problem reading the attributes.
	if r.Then, ok, err = (attributesParser{}).Parse(in); err != nil || !ok {
		err = parse.Error("attribute else if: expected attributes in block, but none were found", in.Position())
		return
	}

	if len(r.Then) == 0 {
		err = parse.Error("attribute else if: invalid content or no attributes were found in the if block", in.Position())
		return
	}

	return r, true, nil
}

var attributeElseExpression parse.Parser[[]Attribute] = attributeElseExpressionParser{}

type attributeElseExpressionParser struct{}

func (attributeElseExpressionParser) Parse(in *parse.Input) (r []Attribute, ok bool, err error) {
	start := in.Index()

	// Strip any initial whitespace.
	_, _, _ = parse.OptionalWhitespace.Parse(in)

	// } else {
	var endElseParser = parse.All(
		parse.Rune('}'),
		parse.OptionalWhitespace,
		parse.String("else"),
		parse.OptionalWhitespace,
		parse.Rune('{'))
	if _, ok, err = endElseParser.Parse(in); err != nil || !ok {
		in.Seek(start)
		return
	}

	// Else contents
	if r, ok, err = (attributesParser{}).Parse(in); err != nil || !ok {
		err = parse.Error("attribute if: expected attributes in else block, but none were found", in.Position())
		return
	}

	return r, true, nil
}
