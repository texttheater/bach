package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type Filter struct {
	Pos           lexer.Position
	FromComponent *FilterFromComponent `"each" @@`
}

func (g *Filter) Ast() (functions.Expression, error) {
	body, err := g.FromComponent.Ast(g.Pos, nil)
	if err != nil {
		return nil, err
	}
	return &functions.MappingExpression{g.Pos, body}, nil
}

type FilterFromComponent struct {
	Pos             lexer.Position
	FromPComponent  *FilterFromPComponent  `( @@`
	FromConditional *FilterFromConditional `| @@ )`
}

func (g *FilterFromComponent) Ast(pos lexer.Position, body functions.Expression) (functions.Expression, error) {
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

func (g *FilterFromPComponent) Ast(pos lexer.Position, body functions.Expression) (functions.Expression, error) {
	component, err := g.PComponent.Ast()
	if err != nil {
		return nil, err
	}
	body = functions.Compose(pos, body, component)
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
	Consequent     *Composition              `"then" ( @@ | "drop" )`
	Pattern        *Pattern                  `( ( ( "elis" @@`
	Guard          *Composition              `      ( "with" @@ )?`
	Condition      *Composition              `    | "elif" @@)`
	FromConsequent *FilterFromConsequentLong `    @@ )`
	Alternative    *Composition              `  | "else" @@ )?`
	FromComponent  *FilterFromComponent      `"ok" ( @@ | "all" )`
}

type FilterFromConsequentShort struct {
	Pos            lexer.Position
	Pattern        *Pattern                   `( ( "elis" @@`
	Guard          *Composition               `    ("with" @@)?`
	Condition      *Composition               `  | "elif" @@ )`
	FromConsequent *FilterFromConsequentShort `  @@`
	FromComponent  *FilterFromComponent       `| ( "ok" @@ | "all" ) )`
}

func (g *FilterFromConditional) Ast(pos lexer.Position, body functions.Expression) (functions.Expression, error) {
	x := &functions.ConditionalExpression{}
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
		x.Pattern = functions.TypePattern{
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
				x.ElisPatterns = append(x.ElisPatterns, pattern)
				if c.Guard != nil {
					guard, err := c.Guard.Ast()
					if err != nil {
						return nil, err
					}
					x.ElisGuards = append(x.ElisGuards, guard)
				} else {
					x.ElisGuards = append(x.ElisGuards, nil)
				}
			} else {
				condition, err := c.Condition.Ast()
				if err != nil {
					return nil, err
				}
				x.ElisPatterns = append(x.ElisPatterns, functions.TypePattern{
					Pos:  c.Pos,
					Type: types.AnyType{},
				})
				x.ElisGuards = append(x.ElisGuards, condition)
			}
			consequent, err = c.FromConsequent.Consequent.Ast()
			if err != nil {
				return nil, err
			}
			x.ElisConsequents = append(x.ElisConsequents, consequent)
			c = c.FromConsequent
		}
		alternative, err := c.Alternative.Ast()
		if err != nil {
			return nil, err
		}
		x.Alternative = alternative
		body = functions.Compose(pos, body, x)
		if c.FromComponent != nil {
			body, err = c.FromComponent.Ast(pos, body)
			if err != nil {
				return nil, err
			}
		}
		return body, nil
	} else { // short form
		x.Consequent = &functions.IdentityExpression{}
		c := g.FromConsequentShort
		for c.FromConsequent != nil { // further elis/elif clauses
			if c.Pattern != nil {
				pattern, err := c.Pattern.Ast()
				if err != nil {
					return nil, err
				}
				x.ElisPatterns = append(x.ElisPatterns, pattern)
				if c.Guard != nil {
					guard, err := c.Guard.Ast()
					if err != nil {
						return nil, err
					}
					x.ElisGuards = append(x.ElisGuards, guard)
				} else {
					x.ElisGuards = append(x.ElisGuards, nil)
				}
			} else {
				condition, err := c.Condition.Ast()
				if err != nil {
					return nil, err
				}
				x.ElisPatterns = append(x.ElisPatterns, functions.TypePattern{
					Pos:  c.Pos,
					Type: types.AnyType{},
				})
				x.ElisGuards = append(x.ElisGuards, condition)
			}
			x.ElisConsequents = append(x.ElisConsequents, &functions.IdentityExpression{pos})
			c = c.FromConsequent
		}
		x.Alternative = &functions.DropExpression{}
		x.UnreachableAlternativeAllowed = true
		body = functions.Compose(pos, body, x)
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
