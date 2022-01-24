package ast

type Source struct {
	ModuleName ModuleName
	FileName   string
	Start      Position
	End        Position
}

type Position struct {
	Line   int
	Column int
}

func MakeSource(
	moduleName ModuleName,
	fileName string,
	start Position,
	end Position,
) *Source {
	return &Source{
		ModuleName: moduleName,
		FileName:   fileName,
		Start:      start,
		End:        end,
	}
}

func MakePosition(
	line int,
	column int,
) Position {
	return Position{
		Line:   line,
		Column: column,
	}
}
