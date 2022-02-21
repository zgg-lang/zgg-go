parser grammar ZggParser;

options {
    tokenVocab=ZggLexer;
    superClass=ZggBaseParser;
}

replItem
    : expr      # replExpr
    | block     # replBlock
    ;

module
    : block
    ;

block
    : (stmt ';'?)*
    ;

codeBlock
    : L_CURLY block R_CURLY
    ;

stmt
    : codeBlock                                 # stmtBlock
    | preIncDec                                 # StmtPreIncDec
    | postIncDec                                # StmtPostIncDec
    | assignExpr                                # stmtAssign
    | callStmt                                  # stmtFuncCall
    | FUNC IDENTIFIER '(' (
          IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?
        | '...' IDENTIFIER ','?
        )? ')'
        codeBlock                               # stmtFuncDefine
    | EXPORT? CLASS className=IDENTIFIER
        ( '(' baseCls+=expr (',' baseCls+=expr)? ')' )?
        L_CURLY memberDef* R_CURLY                       # stmtClassDefine
    | (label=IDENTIFIER ':')? FOR   
            initExpr=expr
            ';'
            checkExpr=expr
            ';'
            nextExpr=expr
            execBlock=codeBlock                 # stmtFor
    | (label=IDENTIFIER ':')? FOR (idIndex=IDENTIFIER ',')? idValue=IDENTIFIER IN
        begin=expr (('..'|'..<') end=expr)?
        execBlock=codeBlock                     # stmtForEach
    | (label=IDENTIFIER ':')? DO    execBlock=codeBlock
        WHILE checkExpr=expr                    # stmtDoWhile
    | (label=IDENTIFIER ':')? WHILE checkExpr=expr execBlock=codeBlock  # stmtWhile
    | CONTINUE (label=IDENTIFIER)?              # stmtContinue
    | BREAK (label=IDENTIFIER)?                 # stmtBreak
    | IF ifCondition codeBlock
        (ELSE IF ifCondition codeBlock)*
        (ELSE codeBlock)?                       # stmtIf
    | SWITCH testValue=expr '{'
        switchCase+
        switchDefault?
        '}'                                     # stmtSwitch
    | RETURN_NONE                               # stmtReturnNone
    | RETURN expr?                              # stmtReturn
    | EXPORT IDENTIFIER                         # stmtExportIdentifier
    | EXPORT IDENTIFIER ':=' expr               # stmtExportExpr
    | EXPORT FUNC IDENTIFIER '('
        ( (IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?)
        | ('...' IDENTIFIER)
        )?
        ')'
        codeBlock                               # stmtExportFuncDefine
    | (DEFER|BLOCK_DEFER) expr '?.'? arguments  # stmtDefer
    | TRY tryBlock=codeBlock (
        CATCH '(' excName=IDENTIFIER ')' catchBlock=codeBlock
        (FINALLY finallyBlock=codeBlock)?
        |
        FINALLY finallyBlock=codeBlock
    )                                           # stmtTry
    | ASSERT expr (',' expr)?                   # stmtAssert
    | EXPORT? EXTEND expr L_CURLY keyValue* R_CURLY     # stmtExtend
    ;

ifCondition
    : (assignExpr ';')? expr
    ;

memberDef
    : STATIC? keyValue
    ;

callStmt
    : expr '?.'? arguments ('??' codeBlock)?
    ;

switchCase
    : CASE whenCondition ':' block FALLTHROUGH?
    ;

switchDefault
    : DEFAULT ':' block
    ;

comparator		: EQUAL | NOT_EQUAL | GTEQ | LTEQ | LT | GT;

expr
    : expr '?.'? arguments                                  	        # exprCall
    | (SINGLE_AT | DOUBLE_AT) IDENTIFIER                                # exprShortImport
    | preIncDec                                             	        # exprPreIncDec
    | postIncDec                                            	        # exprPostIncDec
    | expr '.' field=IDENTIFIER                             	        # exprByField
    | expr '[' index=expr ']'                               	        # exprByIndex
    | IDENTIFIER                                            	        # exprIdentifier
    | literal                                               	        # exprLiteral
    | '-' expr                                              	        # exprNegative
    | '!' expr                                              	        # exprLogicNot
    | '~' expr                                              	        # exprBitNot
    | <assoc=right> expr '**' expr                          	        # exprPow
    | expr op=('*' | '/' | '%') expr                        	        # exprTimesDivMod
    | expr op=('+' | '-') expr                              	        # exprPlusMinus
    | expr op=('<<' | '>>') expr                            	        # exprBitShift
    | expr comparator expr                                              # exprCompare
    | expr '&' expr                                         	        # exprBitAnd
    | expr '|' expr                                         	        # exprBitOr
    | expr '^' expr                                         	        # exprBitXor
    | expr '&&' expr                                        	        # exprLogicAnd
    | expr '||' expr                                        	        # exprLogicOr
    | 'when' L_CURLY
            (expr '->' expr)+
            ('else' '->' expr )?
        R_CURLY                                             	        # exprWhen
    | 'when' expr L_CURLY
            (whenCondition '->' expr)+
            ('else' '->' expr )?
        R_CURLY                                             	        # exprWhenValue
    | condition=expr '?' trueExpr=expr ':' falseExpr=expr   	        # exprQuestion
    | expr '??' expr                                        	        # exprFallback
    | assignExpr                                            	        # exprAssign
    | '(' expr ')'                                          	        # exprSub
    | USE_AT IDENTIFIER expr                                            # exprUseMethod
    | USE_AT codeBlock expr                                             # exprUseBlock
    | USE expr                                              	        # exprUseCloser
    | expr '!'                                              	        # exprAssertError
    ;
    
whenCondition
    : expr (',' expr)*          # whenConditionInList
    | expr ('..'|'..<') expr    # whenConditionInRange
    ;
arguments
    : '(' ( funcArgument ( ',' funcArgument )* ','? )? ')'
    ;

funcArgument
    : ('...'? expr | codeBlock)
    | IDENTIFIER ':' expr
    ;

assignExpr
    : <assoc=right> lval op=('=' | '+=' | '-=' | '*=' | '/=' | '&='
            | '|=' | '^=' | '<<=' | '>>=') expr                     # assignExists
    | <assoc=right> IDENTIFIER ':=' expr                            # assignNew
    | <assoc=right> '[' IDENTIFIER (',' IDENTIFIER )* (',' '...' IDENTIFIER)? ','? ']' ':=' expr
                                                                    # assignNewDeArray
    | <assoc=right> L_CURLY IDENTIFIER (',' IDENTIFIER )* L_CURLY ':=' expr
                                                                    # assignNewDeObject
    | <assoc=right> '...' ':=' expr                                 # assignNewLocal
    ;

preIncDec
    : op=('++' | '--') lval
    ;

postIncDec
    : lval op=('++' | '--')
    ;

lval
    : lval '.' field=IDENTIFIER                             # lvalByField
    | lval '[' index=expr ']'                               # lvalByIndex
    | IDENTIFIER                                            # lvalById
    ;

integer
    :   INT_ZERO        # IntegerZero
    |   INT_DEC         # IntegerDec
    |   INT_HEX         # IntegerHex
    |   INT_OCT         # IntegerOct
    |   INT_BIN         # IntegerBin
    ;

literal
    : integer           # LiteralInteger
    | FLOAT             # LiteralFloat
    | BIGNUM            # LiteralBigNum
    | ('true'|'false')  # LiteralBool
    | stringLiteral     # LiteralString
    | NIL               # LiteralNil
    | UNDEFINED         # LiteralUndefined
    | FUNC '('
        ( (IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?)
        | ('...' IDENTIFIER)
        )? ')' codeBlock               # LiteralFunc            
    | ( '(' (
                IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?
                | '...' IDENTIFIER
            )?
        ')'
      | IDENTIFIER
      ) '=>' expr       # LiteralLambdaExpr
    | ( '(' (
                IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?
                | '...' IDENTIFIER
            )?
        ')'
      | IDENTIFIER
      ) '=>' codeBlock  # LiteralLambdaBlock
    | L_CURLY (objItem (',' objItem)* ','?)? R_CURLY                                      # LiteralObject
    | '[' (arrayItem (',' arrayItem)* ','? )? ']'                                 # LiteralArray
    ;

arrayItem
    : '...'? expr
    ;

objItem
    : keyValue      # ObjItemKV
    | '...' expr    # ObjItemExpanded
    ;

keyValue
    : IDENTIFIER ':' expr    # KVIdKey
    | stringLiteral ':' expr # KVStrKey
    | '[' expr ']' ':' expr  # KVExprKey
    | IDENTIFIER '('
        ( IDENTIFIER (',' IDENTIFIER)* (',' '...' IDENTIFIER)? ','?
        | '...' IDENTIFIER 
        )?
      ')' codeBlock # KVKeyFunc
    | IDENTIFIER             # KVIdOnly
    | '[' expr ']'           # KVExprOnly
    ;

stringLiteral
    : STRING
    | templateString
    ;

templateString
    : QUOTE tsItem* QUOTE
    ;

tsItem
    : TS_RAW                        # tsRaw
    | TS_IDENTIFIER                 # tsIdentifier
    | TS_EXPR_START expr R_CURLY    # tsExpr
    ;
