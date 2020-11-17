<?php 

namespace Ktools\Json;

class Parser {
    
    private int $index;

    private int $len;

    private array $jsonArr;

    private string $json;

    private static array $statusMap = [
        Token::T_BEGIN_DOC  => [
            Token::T_BEGIN_OBJ,
            Token::T_BEGIN_ARR,
            Token::T_STRING,
            Token::T_NUMBER,
            Token::T_BOOLEAN,
            Token::T_NULL,
        ],    
        Token::T_STRING     => ['readString'],
        Token::T_NUMBER     => ['readNumber'],
        Token::T_BOOLEAN    => ['readBoolean'],
        Token::T_NULL       => ['readNull'],
        Token::T_BEGIN_OBJ  => ['readObject'],
        Token::T_END_OBJ    => [
            Token::T_SEP_QUOTE,
            Token::T_END_OBJ,
            Token::T_END_ARR,
        ],
        Token::T_BEGIN_ARR  => ['readArray'],
        Token::T_END_ARR    => [
            Token::T_SEP_QUOTE,
            Token::T_END_OBJ,
            Token::T_END_ARR,
        ],
        Token::T_SEP_COLON  => [
            Token::T_BEGIN_OBJ,
            Token::T_BEGIN_ARR,
            Token::T_SEP_QUOTE,
            Token::T_NUMBER,
            Token::T_BOOLEAN,
            Token::T_NULL,
        ],
        Token::T_SEP_COMMMA => [
            Token::T_SEP_QUOTE
        ],
        Token::T_SEP_QUOTE  => [
            Token::T_STRING
        ],
        Token::T_END_DOC    => ['endParse'],
    ];

    public function nextStatus(int $status) : int {
        $expects = self::$statusMap[$status];
        foreach($expects as $expect) {
            if ($this->hasStatus($expect)) {
                while(true) {
                    if (!$this->hasMore()) {
                        return Token::T_END_DOC; 
                    }
                    $ch = $this->nextChar();
                    if (!$this->isWhiteSpace($ch) || Token::T_SEP_QUOTE == $status) {
                        break; 
                    }
                }
                if ($this->asExpect($ch, $expect)) {
                    return $expect; 
                }
                throw new Error('unexpect character');
            } else {
                return $this->{$expect}();
            }
        }
    }

    public function hasStatus($status) : bool {
        return is_integer($status)
    }

    public function __construct(string $json) {
        $this->index   = 0;
        $this->len     = strlen($json);
        $this->json    = $json; 
        $this->jsonArr = array();
    }

    public function parse() {
        $status = Token::T_BEGIN_DOC;
        while($this->hasMore()) {
            $status = $this->nextStatus($status);
        }
        return $this->jsonArr;
    }

    private function hasMore() : bool {
        return $this->index < $this->len; 
    }

    private function readNumber(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function readTrue() : int {
        $bool = $this->nextFewChar(3);
        if ('rue' === $bool) {
            return 'true';
        }
        throw new Error('Invalid json string');
    }

    private function readFalse() : int {
        $bool = $this->nextFewChar(4);
        if ('alse' === $bool) {
            return 'false';
        }
        throw new Error('Invalid json string');
    }

    private function readNull(&$arr, $index) : string {
        $null = $this->nextFewChar(3);
        if ('ull' === $null) {
            return 'null';
        }
        throw new Error('Invalid json string');
    }

    private function readString() : string {
        $chr = $this->nextChar();
        $index += 1;
        return $index;
    }

    public function nextChar() : string {
        return $this->nextFewChar(1);
    }

    public function peekChar() : string {
        return $this->json[$this->index];
    }

    public function nextFewChar(int $size) : string {
        $this->index += $size;
        return substr($this->json, $this->index, $size);
    } 

    public function asExpect(string $ch, int $status) : bool {
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
