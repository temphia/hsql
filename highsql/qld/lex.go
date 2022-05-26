package qld

import "github.com/araddon/qlbridge/lex"

func sourceMatch(c *lex.Clause, peekWord string, l *lex.Lexer) bool {
	switch peekWord {
	case "(":
		return true
	case "select":
		return true
	case "left", "right", "inner", "outer", "join":
		return true
	}
	return false
}

var (
	SqlSelect = []*lex.Clause{
		{Token: lex.TokenSelect, Lexer: lex.LexSelectClause, Name: "sqlSelect.Select"},
		{Token: lex.TokenFrom, Lexer: lex.LexTableReferenceFirst, Optional: true, Repeat: false, Clauses: fromSource, Name: "sqlSelect.From"},
		{Token: lex.TokenWhere, Lexer: lex.LexConditionalClause, Optional: true, Clauses: whereQuery, Name: "sqlSelect.where"},
		{KeywordMatcher: sourceMatch, Optional: true, Repeat: true, Clauses: moreSources, Name: "sqlSelect.sources"},
		{Token: lex.TokenGroupBy, Lexer: lex.LexColumns, Optional: true, Name: "sqlSelect.groupby"},
		{Token: lex.TokenHaving, Lexer: lex.LexConditionalClause, Optional: true, Name: "sqlSelect.having"},
		{Token: lex.TokenOrderBy, Lexer: lex.LexOrderByColumn, Optional: true, Name: "sqlSelect.orderby"},
		{Token: lex.TokenLimit, Lexer: lex.LexLimit, Optional: true, Name: "sqlSelect.limit"},
		{Token: lex.TokenOffset, Lexer: lex.LexNumber, Optional: true, Name: "sqlSelect.offset"},
		{Token: lex.TokenWith, Lexer: lex.LexJsonOrKeyValue, Optional: true, Name: "sqlSelect.with"},
		{Token: lex.TokenAlias, Lexer: lex.LexIdentifier, Optional: true, Name: "sqlSelect.alias"},
		{Token: lex.TokenEOF, Lexer: lex.LexEndOfStatement, Optional: false, Name: "sqlSelect.eos"},
	}
	fromSource = []*lex.Clause{
		{KeywordMatcher: sourceMatch, Lexer: lex.LexTableReferenceFirst, Name: "fromSource.matcher"},
		{Token: lex.TokenSelect, Lexer: lex.LexSelectClause, Name: "fromSource.Select"},
		{Token: lex.TokenFrom, Lexer: lex.LexTableReferenceFirst, Optional: true, Repeat: true, Name: "fromSource.From"},
		{Token: lex.TokenWhere, Lexer: lex.LexConditionalClause, Optional: true, Name: "fromSource.Where"},
		{Token: lex.TokenHaving, Lexer: lex.LexConditionalClause, Optional: true, Name: "fromSource.having"},
		{Token: lex.TokenGroupBy, Lexer: lex.LexColumns, Optional: true, Name: "fromSource.GroupBy"},
		{Token: lex.TokenOrderBy, Lexer: lex.LexOrderByColumn, Optional: true, Name: "fromSource.OrderBy"},
		{Token: lex.TokenLimit, Lexer: lex.LexLimit, Optional: true, Name: "fromSource.Limit"},
		{Token: lex.TokenOffset, Lexer: lex.LexNumber, Optional: true, Name: "fromSource.Offset"},
		{Token: lex.TokenRightParenthesis, Lexer: lex.LexEndOfSubStatement, Optional: true, Name: "fromSource.EndParen"},
		{Token: lex.TokenAs, Lexer: lex.LexIdentifier, Optional: true, Name: "fromSource.As"},
		{Token: lex.TokenOn, Lexer: lex.LexConditionalClause, Optional: true, Name: "fromSource.On"},
	}

	moreSources = []*lex.Clause{
		{KeywordMatcher: sourceMatch, Lexer: lex.LexJoinEntry, Name: "moreSources.JoinEntry"},
		{Token: lex.TokenSelect, Lexer: lex.LexSelectClause, Optional: true, Name: "moreSources.Select"},
		{Token: lex.TokenFrom, Lexer: lex.LexTableReferenceFirst, Optional: true, Repeat: true, Name: "moreSources.From"},
		{Token: lex.TokenWhere, Lexer: lex.LexConditionalClause, Optional: true, Name: "moreSources.Where"},
		{Token: lex.TokenHaving, Lexer: lex.LexConditionalClause, Optional: true, Name: "moreSources.Having"},
		{Token: lex.TokenGroupBy, Lexer: lex.LexColumns, Optional: true, Name: "moreSources.GroupBy"},
		{Token: lex.TokenOrderBy, Lexer: lex.LexOrderByColumn, Optional: true, Name: "moreSources.OrderBy"},
		{Token: lex.TokenLimit, Lexer: lex.LexLimit, Optional: true, Name: "moreSources.Limit"},
		{Token: lex.TokenOffset, Lexer: lex.LexNumber, Optional: true, Name: "moreSources.Offset"},
		{Token: lex.TokenRightParenthesis, Lexer: lex.LexEndOfSubStatement, Optional: false, Name: "moreSources.EndParen"},
		{Token: lex.TokenAs, Lexer: lex.LexIdentifier, Optional: true, Name: "moreSources.As"},
		{Token: lex.TokenOn, Lexer: lex.LexConditionalClause, Optional: true, Name: "moreSources.On"},
	}

	whereQuery = []*lex.Clause{
		{Token: lex.TokenSelect, Lexer: lex.LexSelectClause, Name: "whereQuery.Select"},
		{Token: lex.TokenFrom, Lexer: lex.LexTableReferences, Optional: true, Repeat: true, Name: "whereQuery.From"},
		{Token: lex.TokenWhere, Lexer: lex.LexConditionalClause, Optional: true, Name: "whereQuery.Where"},
		{Token: lex.TokenHaving, Lexer: lex.LexConditionalClause, Optional: true, Name: "whereQuery.Having"},
		{Token: lex.TokenGroupBy, Lexer: lex.LexColumns, Optional: true, Name: "whereQuery.GroupBy"},
		{Token: lex.TokenOrderBy, Lexer: lex.LexOrderByColumn, Optional: true, Name: "whereQuery.OrderBy"},
		{Token: lex.TokenLimit, Lexer: lex.LexNumber, Optional: true, Name: "whereQuery.Limit"},
		{Token: lex.TokenRightParenthesis, Lexer: lex.LexEndOfSubStatement, Optional: false, Name: "whereQuery.EOS"},
	}

	partialSql = &lex.Dialect{
		Statements: []*lex.Clause{
			{Token: lex.TokenSelect, Clauses: SqlSelect},
			{Token: lex.TokenFrom, Clauses: SqlSelect},
		},
	}
)
