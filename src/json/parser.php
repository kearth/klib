<?php 

namespace Ktools\Json;

class Parser {
    
    private int $index;

    private array $jsonArr;

    private string $json;

    public function __construct(string $json) {
        $this->index   = 0;
        $this->json    = $json; 
        $this->jsonArr = array();
    }

    public function parse() {
        $len = strlen($this->json);
        $lexer = new Lexer($this->jsonArr, $this->json);
        while($this->hasMore($len)) {
            $this->index = $lexer->getToken($this->index);
        }
        return $this->jsonArr;
    }

    private function hasMore(int $len) : bool {
        return $this->index < $len; 
    }
}
