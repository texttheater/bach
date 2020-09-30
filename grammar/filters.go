package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

type Filter struct {
	Pos           lexer.Position
	FromComponent *FilterFromComponent `"each" @@`
}

func (g *Filter) Ast() (expressions.Expression, error) {
	body, err := g.FromComponent.Ast(g.Pos, nil)
	if err != nil {
		return nil, err
	}
	return &expressions.MappingExpression{g.Pos, body}, nil
}

type FilterFromComponent struct {
	Pos             lexer.Position
	FromPComponent  *FilterFromPComponent  `( @@`
	FromConditional *FilterFromConditional `| @@ )`
}

func (g *FilterFromComponent) Ast(pos lexer.Position, body expressions.Expression) (expressions.Expression, error) {
	if g.FromPComponent != nil {
		return g.FromPComponent.Ast(pos, body)
	}
	if g.FromConditional != nil {
		return g.FromConditional.Ast(pos, body)
	}
	panic("invalid component")
}

type FilterFromPComponent struct {
	Pos           lexer.Position
	PComponent    *PComponent          `@@`
	FromComponent *FilterFromComponent `( @@ | "all" )`
}

func (g *FilterFromPComponent) Ast(pos lexer.Position, body expressions.Expression) (expressions.Expression, error) {
	component, err := g.PComponent.Ast()
	if err != nil {
		return nil, err
	}
	body = expressions.Compose(pos, body, component)
	if g.FromComponent != nil {
		body, err = g.FromComponent.Ast(pos, body)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

type FilterFromConditional struct {
	Pos                 lexer.Position
	Pattern             *Pattern                   `( "is" @@`
	Guard               *Composition               `  ( "with" @@ )?`
	Condition           *Composition               `| "if" @@ )`
	FromConsequentLong  *FilterFromConsequentLong  `( @@`
	FromConsequentShort *FilterFromConsequentShort `| @@ )`
}

type FilterFromConsequentLong struct {
	Pos            lexer.Position
	Consequent     *QComposition             `"then" @@`
	Pattern        *Pattern                  `( ( "elis" @@`
	Guard          *Composition              `    ( "with" @@ )?`
	Condition      *Composition              `  | "elif" @@ )`
	FromConsequent *FilterFromConsequentLong `  @@`
	Alternative    *QComposition             `| ( "else" @@ )?`
	FromComponent  *FilterFromComponent      `  "ok" ( @@ | "all" ) )`
}

type FilterFromConsequentShort struct {
	Pos            lexer.Position
	Pattern        *Pattern                   `( ( "elis" @@`
	Guard          *Composition               `    ("with" @@)?`
	Condition      *Composition               `  | "elif" @@ )`
	FromConsequent *FilterFromConsequentShort `  @@`
	FromComponent  *FilterFromComponent       `| ( "ok" @@ | "all" ) )`
}

func (g *FilterFromConditional) Ast(pos lexer.Position, body expressions.Expression) (expressions.Expression, error) {
	x := &expressions.ConditionalExpression{}
	x.Pos = g.Pos
	if g.Pattern != nil {
		pattern, err := g.Pattern.Ast()
		if err != nil {
			return nil, err
		}
		x.Pattern = pattern
		if g.Guard != nil {
			guard, err := g.Guard.Ast()
			if err != nil {
				return nil, err
			}
			x.Guard = guard
		}
	} else {
		condition, err := g.Condition.Ast()
		if err != nil {
			return nil, err
		}
		x.Pattern = expressions.TypePattern{
			Pos:  g.Pos,
			Type: types.AnyType{},
		}
		x.Guard = condition
	}
	if g.FromConsequentLong != nil { // long form
		consequent, err := g.FromConsequentLong.Consequent.Ast()
		if err != nil {
			return nil, err
		}
		x.Consequent = consequent
		c := g.FromConsequentLong
		for c.FromConsequent != nil { // further elis/elif clauses
			if c.Pattern != nil {
				pattern, err := c.Pattern.Ast()
				if err != nil {
					return nil, err
				}
				x.AlternativePatterns = append(x.AlternativePatterns, pattern)
				if c.Guard != nil {
					guard, err := c.Guard.Ast()
					if err != nil {
						return nil, err
					}
					x.AlternativeGuards = append(x.AlternativeGuards, guard)
				} else {
					x.AlternativeGuards = append(x.AlternativeGuards, nil)
				}
			} else {
				condition, err := c.Condition.Ast()
				if err != nil {
					return nil, err
				}
				x.AlternativePatterns = append(x.AlternativePatterns, expressions.TypePattern{
					Pos:  c.Pos,
					Type: types.AnyType{},
				})
				x.AlternativeGuards = append(x.AlternativeGuards, condition)
			}
			consequent, err = c.FromConsequent.Consequent.Ast()
			if err != nil {
				return nil, err
			}
			x.AlternativeConsequents = append(x.AlternativeConsequents, consequent)
			c = c.FromConsequent
		}
		if c.Alternative != nil {
			alternative, err := c.Alternative.Ast()
			if err != nil {
				return nil, err
			}
			x.Alternative = alternative
		}
		body = expressions.Compose(pos, body, x)
		if c.FromComponent != nil {
			body, err = c.FromComponent.Ast(pos, body)
			if err != nil {
				return nil, err
			}
		}
		return body, nil
	} else { // short form
		x.Consequent = &expressions.IdentityExpression{}
		c := g.FromConsequentShort
		for c.FromConsequent != nil { // further elis/elif clauses
			if c.Pattern != nil {
				pattern, err := c.Pattern.Ast()
				if err != nil {
					return nil, err
				}
				x.AlternativePatterns = append(x.AlternativePatterns, pattern)
				if c.Guard != nil {
					guard, err := c.Guard.Ast()
					if err != nil {
						return nil, err
					}
					x.AlternativeGuards = append(x.AlternativeGuards, guard)
				} else {
					x.AlternativeGuards = append(x.AlternativeGuards, nil)
				}
			} else {
				condition, err := c.Condition.Ast()
				if err != nil {
					return nil, err
				}
				x.AlternativePatterns = append(x.AlternativePatterns, expressions.TypePattern{
					Pos:  c.Pos,
					Type: types.AnyType{},
				})
				x.AlternativeGuards = append(x.AlternativeGuards, condition)
			}
			x.AlternativeConsequents = append(x.AlternativeConsequents, &expressions.IdentityExpression{pos})
			c = c.FromConsequent
		}
		x.Alternative = &expressions.DropExpression{}
		x.UnreachableAlternativeAllowed = true
		body = expressions.Compose(pos, body, x)
		if c.FromComponent != nil {
			var err error
			body, err = c.FromComponent.Ast(pos, body)
			if err != nil {
				return nil, err
			}
		}
		return body, nil
	}
}
