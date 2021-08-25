# Grammar

```
start -> "package" IDENTIFIER line-sep import* line-sep decl-stmt*;
import -> "import" IDENTIFIER ("." IDENTIFIER)* ("{" (IDENTIFIER ("," IDENTIFIER)*)? "}")?;

line-sep -> "\n" | ";";

data-decl -> "data" IDENTIFIER "{" (data-decl-property (line-sep data-decl-property)*)? "}";
data-decl-property -> IDENTIFIER (IDENTIFIER ("," IDENTIFIER)*)?;
enum-decl -> "enum" IDENTIFIER "{" enum-decl-case* "}";
enum-decl-case -> data-decl | enum-decl | IDENTIFIER;

func-decl -> "func" IDENTIFIER ("=>" decl-stmt | func-expr);
let-decl -> "let" IDENTIFIER = expr;

func-expr -> "{" func-expr-args "=>" decl-stmt* "}";
func-expr-args -> IDENTIFIER ("," IDENTIFIER)*;

decl-stmt -> data-decl |
    enum-decl |
    func-decl |
    let-decl |
    expr;

expr -> simple-expr (simple-expr ("," simple-expr)*)?;

simple-expr -> literal ("." IDENTIFIER)* simple-expr?;
literal -> IDENTIFIER |
    NUMBER |
    STRING |
    array |
    grouping |
    func-expr;

grouping -> "(" expr ")";
array -> "[" (simple-expr ("," simple-expr))? "]";

// TODO: all operators missing!
// TODO: de-/constructors missing!

line-sep -> "\n" | ";";

```

## Precedence

1. `.` member access
2. `-`, `!` unary negation
3. `*`, `/` multiplicative
4. `+`, `-` additive
5. `==`, `!=`, `>`, `>=`, `<=`, `<` comparisions
6. `,` argument separation
