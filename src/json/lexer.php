<?php 

namespace Ktools\Json;

class Lexer {

    private int $index;

    private int $len;

    private string $json;

    private int $status;

    private int $initStatus = 0; 

    public function __construct(string $json) {
        $this->json    = $json; 
        $this->len     = strlen($json);
        $this->status  = 0;
    }    

    public function isWhiteSpace(string $ch) : bool {
        return ' ' === $ch; 
    }

    public function nextToken() : object {
        while(true) {
            if (!$this->hasMore()) {
                return new Token(Token.T_END_DOC, null); 
            }
            $ch = $this->nextChar();
            if(!isWhiteSpace($ch)) {
                break; 
            }
        }
        switch ($ch) {
        case "{":       
            return new Token(Token.T_BEGIN_OBJ, null);
        case "}":       
            return new Token(Token.T_END_OBJ, null);
        case "[":        
            return new Token(Token.T_BEGIN_ARR, null);
        case "]":        
            return new Token(Token.T_END_ARR, null);
        case ',':
            return new Token(Token.T_SEP_COMMMA, null);
        case ':':
            return new Token(Token.T_SEP_COLON, null);
        case 'n':
            return new Token(Token.T_NULL, $this->readNull());
        case 't':
            return new Token(Token.T_BOOLEAN, $this->readTrue());
        case 'f':
            return new Token(Token.T_BOOLEAN, $this->readFalse());
        case '"':
            return new Token(Token.T_STRING, $this->readString());
        default:
            if (is_numeric($ch)) {
                return new Token(Token.T_NUMBER, $this->readNumber());
            }
        }
        throw new Error('Illegal character');
    }
}
